package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Book struct {
	Id   int64
	Text string
}

var items = []Book{
	{Id: 1, Text: "Ultima noapte de dragoste, întâia noapte de război de Camil Petrescu"},
	{Id: 2, Text: "Pădurea spânzuraţilor de Liviu Rebreanu"},
	{Id: 3, Text: "Moromeții I de Marin Preda"},
	{Id: 4, Text: "Maitreyi de Mircea Eliade"},
	{Id: 5, Text: "Enigma Otiliei de George Călinescu"},
	{Id: 6, Text: "La țigănci de Mircea Eliade"},
	{Id: 7, Text: "Moara cu noroc de Ioan Slavici"},
	{Id: 8, Text: "Amintiri din copilărie de Ion Creangă"},
	{Id: 9, Text: "Patul lui Procust de Camil Petrescu"},
	{Id: 10, Text: "Elevul Dima dintr-a VII-A de Mihail Drumeș"},
	{Id: 11, Text: "Întoarcerea din rai de Mircea Eliade"},
	{Id: 12, Text: "La hanul lui Mânjoală de Ion Luca Caragiale"},
	{Id: 13, Text: "O scrisoare pierdută de Ion Luca Caragiale"},
	{Id: 14, Text: "Ion de Liviu Rebreanu"},
	{Id: 15, Text: "Baltagul de Mihail Sadoveanu"},
}

func Test_Engine(t *testing.T) {
	engine := NewEngine()
	engine.SetTolerance(2)
	for _, item := range items {
		engine.SetItem(item.Id, item.Text)
	}
	tests := []struct {
		query    string
		expected SearchResult
	}{
		{"maitreyi", SearchResult{Items: []int64{4}, Total: 1, Pages: 1}},
		{"eliade", SearchResult{Items: []int64{11, 6, 4}, Total: 3, Pages: 1}},
		{"Patul lui", SearchResult{Items: []int64{9, 12, 11, 10, 7}, Total: 8, Pages: 2}},
		{"spânzuraţilor", SearchResult{Items: []int64{2}, Total: 1, Pages: 1}},
		{"amintiri din copilărie", SearchResult{Items: []int64{8, 11, 10, 5, 15}, Total: 15, Pages: 3}},
		{"xyz zyx", SearchResult{Items: []int64{}, Total: 0, Pages: 0}},
		{"din", SearchResult{Items: []int64{11, 8, 15, 14, 13}, Total: 15, Pages: 3}},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			actual := engine.Search(SearchOptions{Query: test.query, Limit: 5})
			assert.Equal(t, test.expected, actual)
		})
	}

	t.Run("offset", func(t *testing.T) {
		actual := engine.Search(SearchOptions{
			Query:  "de",
			Limit:  5,
			Offset: 5,
			Ignore: []int64{15},
		})
		assert.Equal(t, SearchResult{Items: []int64{9, 8, 7, 6, 5}, Total: 14, Pages: 3}, actual)
	})

	t.Run("Ignore ids", func(t *testing.T) {
		actual := engine.Search(SearchOptions{
			Query:  "maitreyi",
			Limit:  5,
			Ignore: []int64{4},
		})
		assert.Equal(t, SearchResult{Items: []int64{}}, actual)
	})

	engine.SetItem(16, "Ciocoii vechi și noi de Nicolae Filimon")
	t.Run("SetItem", func(t *testing.T) {
		actual := engine.Search(SearchOptions{
			Query: "Ciocoii vechi",
			Limit: 5,
		})
		assert.Equal(t, SearchResult{Items: []int64{16}, Total: 1, Pages: 1}, actual)
	})

	engine.DeleteItem(7)
	t.Run("DeleteItem", func(t *testing.T) {
		actual := engine.Search(SearchOptions{
			Query: "Moara",
			Limit: 5,
		})
		assert.Equal(t, SearchResult{Items: []int64{}}, actual)
	})
}

func Test_Tokenize(t *testing.T) {
	tokens := Tokenize("Țară, România, școală, mâine! 123.4 ăĂâÂîÎșşȘŞțţȚŢ")
	expected := [][]rune{
		{'t', 'a', 'r', 'a'},
		{'r', 'o', 'm', 'a', 'n', 'i', 'a'},
		{'s', 'c', 'o', 'a', 'l', 'a'},
		{'m', 'a', 'i', 'n', 'e'},
		{'1', '2', '3'},
		{'4'},
		{'a', 'a', 'a', 'a', 'i', 'i', 's', 's', 's', 's', 't', 't', 't', 't'},
	}
	assert.Equal(t, expected, tokens)
}

func Test_LevenshteinDistance(t *testing.T) {
	tests := []struct {
		a        []rune
		b        []rune
		expected int
	}{
		{[]rune(""), []rune(""), 0},
		{[]rune("a"), []rune(""), 1},
		{[]rune(""), []rune("a"), 1},
		{[]rune("a"), []rune("a"), 0},
		{[]rune("a"), []rune("ab"), 1},
		{[]rune("ab"), []rune("a"), 1},
		{[]rune("sitting"), []rune("kitten"), 3},
		{[]rune("kitten"), []rune("sitting"), 3},
		{[]rune("sitting"), []rune("sitting"), 0},
		{[]rune("lawn"), []rune("flaw"), 2},
		{[]rune("flaw"), []rune("lawn"), 2},
		{[]rune("kit"), []rune("kitten"), 3},
		{[]rune("kitten"), []rune("kit"), 3},
	}
	for _, test := range tests {
		t.Run("LevenshteinDistance", func(t *testing.T) {
			assert.Equal(t, test.expected, LevenshteinDistance(test.a, test.b))
		})
	}
}
