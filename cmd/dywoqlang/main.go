package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dywoq/lang/pkg/parser"
	"github.com/dywoq/lang/pkg/scanner"
	"github.com/dywoq/lang/pkg/symbol"
)

func main() {
	f, err := os.Open("./main.dl")
	if err != nil {
		panic(err)
	}

	s, err := scanner.New(f)
	if err != nil {
		panic(err)
	}
	tokens, err := s.Scan("main.dl")
	if err != nil {
		panic(err)
	}

	p := parser.New(tokens)
	if err != nil {
		panic(err)
	}

	tree, err := p.Parse("main.dl")
	if err != nil {
		panic(err)
	}

	i := symbol.NewInfo()

	symbols, err := i.Collect(tree)
	if err != nil {
		panic(err)
	}

	for key, symbol := range symbols {
		content, err := json.MarshalIndent(symbol, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s - %s\n", key, string(content))
	}
}
