# Search

In-memory fuzzy search engine for Golang.

[![Go Reference](https://pkg.go.dev/badge/github.com/radulucut/search.svg)](https://pkg.go.dev/github.com/radulucut/search)
![Test](https://github.com/radulucut/search/actions/workflows/test.yml/badge.svg)

## Installation

```bash
go get github.com/radulucut/search
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/radulucut/search"
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

func main() {
	// Create a new search engine
	engine := search.NewEngine(items, func(item Book) (int64, string) {
		return item.Id, item.Text
	}, nil)
	engine.SetTolerance(2)

	// Search for a book
	results := engine.Search("Eliade", 5)

	// Print the results
	// Output: []int64{11, 6, 4}
	fmt.Println(results)
}

```

## TODO

- [ ] Add stemming support
- [ ] Add tf-idf support
