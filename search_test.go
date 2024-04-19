package search

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		expected []int64
	}{
		{"maitreyi", []int64{4}},
		{"eliade", []int64{11, 6, 4}},
		{"Patul lui", []int64{9, 12, 11, 10, 7}},
		{"spânzuraţilor", []int64{2}},
		{"amintiri din copilărie", []int64{8, 11, 10, 5, 15}},
		{"xyz zyx", []int64{}},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			actual := engine.Search(test.query, 5, nil)
			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("Ignore ids", func(t *testing.T) {
		actual := engine.Search("maitreyi", 5, []int64{4})
		expected := []int64{}
		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	engine.SetItem(16, "Ciocoii vechi și noi de Nicolae Filimon")
	t.Run("SetItem", func(t *testing.T) {
		actual := engine.Search("Ciocoii vechi", 5, nil)
		expected := []int64{16}
		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	engine.DeleteItem(7)
	t.Run("DeleteItem", func(t *testing.T) {
		actual := engine.Search("Moara", 5, nil)
		expected := []int64{}
		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
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
	if diff := cmp.Diff(expected, tokens); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
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
			if diff := cmp.Diff(test.expected, LevenshteinDistance(test.a, test.b)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
