package search

import (
	"math"
	"slices"
	"sync"
)

type Engine struct {
	sync.RWMutex
	tolerance int
	items     []item
}

type item struct {
	Id   int64
	Text [][]rune
}

func NewEngine() *Engine {
	engine := &Engine{
		tolerance: 1,
	}
	return engine
}

// Set the maximum number of typos per word allowed.
// The default value is 1.
func (e *Engine) SetTolerance(tolerance int) {
	e.Lock()
	defer e.Unlock()
	e.tolerance = tolerance
}

// Add a new item to the search engine.
func (e *Engine) SetItem(id int64, text string) {
	e.Lock()
	defer e.Unlock()
	index, ok := binarySearch(e.items, id)
	if ok {
		e.items[index].Text = tokenize(text)
		return
	}
	e.items = append(e.items, item{})
	copy(e.items[index+1:], e.items[index:])
	e.items[index] = item{
		Id:   id,
		Text: tokenize(text),
	}
}

// Remove an item from the search engine.
func (e *Engine) DeleteItem(id int64) {
	e.Lock()
	defer e.Unlock()
	index, ok := binarySearch(e.items, id)
	if !ok {
		return
	}
	e.items = append(e.items[:index], e.items[index+1:]...)
}

type itemScore struct {
	id    int64
	score int
}

type SearchOptions struct {
	Query  string
	Limit  int
	Offset int
	Ignore []int64
}

type SearchResult struct {
	Items []int64
	Total int
	Pages int
}

// Search finds the most similar items to the given query.
// limit is the maximum number of items to return.
// ignore is a list of item ids to ignore.
func (e *Engine) Search(opts SearchOptions) SearchResult {
	var ignoreMap map[int64]struct{}
	hasIgnore := false
	if len(opts.Ignore) != 0 {
		hasIgnore = true
		ignoreMap = make(map[int64]struct{})
		for i := range opts.Ignore {
			ignoreMap[opts.Ignore[i]] = struct{}{}
		}
	}
	q := tokenize(opts.Query)
	e.RLock()
	defer e.RUnlock()
	scores := make([]*itemScore, 0)
	for i := range e.items {
		id := e.items[i].Id
		if hasIgnore {
			if _, ok := ignoreMap[id]; ok {
				continue
			}
		}
		score := e.score(q, e.items[i].Text)
		if score == -1 {
			continue
		}
		scores = append(scores, &itemScore{id: id, score: score})
	}
	slices.SortFunc(scores, func(a, b *itemScore) int {
		if a.score < b.score {
			return -1
		}
		if a.score > b.score {
			return 1
		}
		if a.id > b.id {
			return -1
		}
		if a.id < b.id {
			return 1
		}
		return 0
	})
	l := len(scores)
	limit := min(opts.Offset+opts.Limit, l)
	res := make([]int64, 0, limit)
	for i := opts.Offset; i < limit; i++ {
		res = append(res, scores[i].id)
	}
	pages := l / opts.Limit
	if l%opts.Limit != 0 {
		pages++
	}
	return SearchResult{
		Items: res,
		Total: l,
		Pages: pages,
	}
}

func (e *Engine) score(q, b [][]rune) int {
	var score int
	skip := true
	for i := range q {
		best := math.MaxInt
		for j := range b {
			best = min(best, levenshteinDistance(q[i], b[j]))
		}
		if best <= e.tolerance {
			skip = false
		}
		score += best
	}
	if skip {
		return -1
	}
	return score
}

func binarySearch(s []item, target int64) (int, bool) {
	n := len(s)
	l, r := 0, n
	for l < r {
		m := int(uint(l+r) >> 1)
		if s[m].Id < target {
			l = m + 1
		} else {
			r = m
		}
	}
	return l, l < n && s[l].Id == target
}
