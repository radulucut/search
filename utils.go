package search

import "unicode"

type TokenizeFunc func(input string) [][]rune

func Tokenize(input string) [][]rune {
	var tokens [][]rune
	var token []rune
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			switch r {
			// Normalize Romanian diacritics.
			case 'ă', 'Ă':
				r = 'a'
			case 'â', 'Â':
				r = 'a'
			case 'î', 'Î':
				r = 'i'
			case 'ș', 'ş', 'Ș', 'Ş':
				r = 's'
			case 'ț', 'ţ', 'Ț', 'Ţ':
				r = 't'
			}
			token = append(token, unicode.ToLower(r))
		} else {
			if len(token) > 0 {
				tokens = append(tokens, token)
				token = nil
			}
		}
	}
	if len(token) > 0 {
		tokens = append(tokens, token)
	}
	return tokens
}

func LevenshteinDistance(a, b []rune) int {
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
