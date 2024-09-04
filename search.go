package search

import (
	"math"
	"slices"
	"sync"
)

type Engine struct {
	sync.RWMutex
	items     map[int64][][]rune
	tokenize  TokenizeFunc
	tolerance int
}

func NewEngine() *Engine {
	engine := &Engine{
		items:     make(map[int64][][]rune),
		tolerance: 1,
		tokenize:  Tokenize,
	}
	return engine
}

// Set custom tokenize function.
func (e *Engine) SetTokenizeFunc(f TokenizeFunc) {
	e.Lock()
	defer e.Unlock()
	e.tokenize = f
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
	e.items[id] = e.tokenize(text)
}

// Remove an item from the search engine.
func (e *Engine) DeleteItem(id int64) {
	e.Lock()
	defer e.Unlock()
	delete(e.items, id)
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

// Search finds the most similar items to the given query.
// limit is the maximum number of items to return.
// ignore is a list of item ids to ignore.
func (e *Engine) Search(opts SearchOptions) []int64 {
	var ignoreMap map[int64]struct{}
	hasIgnore := false
	if len(opts.Ignore) != 0 {
		hasIgnore = true
		ignoreMap = make(map[int64]struct{})
		for i := range opts.Ignore {
			ignoreMap[opts.Ignore[i]] = struct{}{}
		}
	}
	q := e.tokenize(opts.Query)
	e.RLock()
	defer e.RUnlock()
	scores := make([]*itemScore, 0)
	for id := range e.items {
		if hasIgnore {
			if _, ok := ignoreMap[id]; ok {
				continue
			}
		}
		score := e.score(q, e.items[id])
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
	limit := min(opts.Offset+opts.Limit, len(scores))
	res := make([]int64, 0, limit)
	for i := opts.Offset; i < limit; i++ {
		res = append(res, scores[i].id)
	}
	return res
}

func (e *Engine) score(q, b [][]rune) int {
	var score int
	skip := true
	for i := range q {
		best := math.MaxInt
		for j := range b {
			best = min(best, LevenshteinDistance(q[i], b[j]))
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
