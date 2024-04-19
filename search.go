package search

import (
	"sort"
	"sync"
)

type Engine struct {
	items     map[int64][][]rune
	tokenize  TokenizeFunc
	tolerance int
	mx        sync.RWMutex
}

func NewEngine() *Engine {
	engine := &Engine{
		items:     make(map[int64][][]rune),
		tolerance: 1,
		tokenize:  Tokenize,
		mx:        sync.RWMutex{},
	}
	return engine
}

// Set custom tokenize function.
func (e *Engine) SetTokenizeFunc(f TokenizeFunc) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.tokenize = f
}

// Set the maximum number of typos per word allowed.
// The default value is 1.
func (e *Engine) SetTolerance(tolerance int) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.tolerance = tolerance
}

// Add a new item to the search engine.
func (e *Engine) SetItem(id int64, text string) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.items[id] = e.tokenize(text)
}

// Remove an item from the search engine.
func (e *Engine) DeleteItem(id int64) {
	e.mx.Lock()
	defer e.mx.Unlock()
	delete(e.items, id)
}

type itemScore struct {
	id    int64
	score int
}

// Search finds the most similar items to the given query.
// limit is the maximum number of items to return.
// ignore is a list of item ids to ignore.
func (e *Engine) Search(query string, limit int, ignore []int64) []int64 {
	var ignoreMap map[int64]struct{}
	hasIgnore := false
	if len(ignore) != 0 {
		hasIgnore = true
		ignoreMap = make(map[int64]struct{})
		for i := range ignore {
			ignoreMap[ignore[i]] = struct{}{}
		}
	}

	q := e.tokenize(query)
	e.mx.RLock()
	defer e.mx.RUnlock()
	scores := make([]itemScore, 0)
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
		scores = append(scores, itemScore{id: id, score: score})
	}
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].score == scores[j].score {
			return scores[i].id > scores[j].id
		} else {
			return scores[i].score < scores[j].score
		}
	})
	limit = min(limit, len(scores))
	res := make([]int64, 0, limit)
	for i := 0; i < limit; i++ {
		res = append(res, scores[i].id)
	}
	return res
}

func (e *Engine) score(q, b [][]rune) int {
	var score int
	skip := true
	for i := range q {
		best := (1<<63 - 1)
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
