package search

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func tokenize(input string) [][]rune {
	var tokens [][]rune
	var token []rune
	for _, r := range normalize(input) {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if len(token) > 0 {
				tokens = append(tokens, token)
				token = nil
			}
			continue
		}
		token = append(token, unicode.ToLower(r))
	}
	if len(token) > 0 {
		tokens = append(tokens, token)
	}
	return tokens
}

func normalize(s string) string {
	r, _, err := transform.String(transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	), s)
	if err != nil {
		return s
	}
	return r
}

func levenshteinDistance(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	if len(a) > len(b) {
		a, b = b, a
	}
	la, lb := len(a), len(b)
	row := make([]int, la+1)
	for i := 1; i <= la; i++ {
		row[i] = i
	}
	for i := 1; i <= lb; i++ {
		prev := i
		for j := 1; j <= la; j++ {
			curr := row[j-1]
			if b[i-1] != a[j-1] {
				curr = min(row[j-1]+1, prev+1, row[j]+1)
			}
			row[j-1] = prev
			prev = curr
		}
		row[la] = prev
	}
	return row[la]
}
