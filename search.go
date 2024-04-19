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

// NewEngine creates a new search engine.
//
// items is a slice of items to be indexed.
//
// mapFunc is a function that maps an item to an id and a string.
//
// tokenizeFunc is an optional function that tokenizes a string into words.
func NewEngine[T any](
	items []T,
	mapFunc func(T) (int64, string),
	tokenizeFunc TokenizeFunc,
) *Engine {
	engine := &Engine{
		items:     make(map[int64][][]rune),
		tolerance: 1,
		tokenize:  Tokenize,
		mx:        sync.RWMutex{},
	}
	if tokenizeFunc != nil {
		engine.tokenize = tokenizeFunc
	}
	for i := range items {
		id, s := mapFunc(items[i])
		engine.items[id] = engine.tokenize(s)
	}
	return engine
}

// SetTolerance sets the maximum number of typos per word allowed.
// The default value is 1.
func (e *Engine) SetTolerance(tolerance int) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.tolerance = tolerance
}

// SetItem adds a new item to the search engine.
func (e *Engine) SetItem(id int64, text string) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.items[id] = e.tokenize(text)
}

// DeleteItem removes an item from the search engine.
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
func (e *Engine) Search(query string, limit int) []int64 {
	q := e.tokenize(query)
	e.mx.RLock()
	defer e.mx.RUnlock()
	scores := make([]itemScore, 0)
	for id := range e.items {
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
