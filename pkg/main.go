package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dywoq/lang/pkg/parser"
	"github.com/dywoq/lang/pkg/scanner"
)

func main() {
	f, err := os.Open("./main.dl")
	if err != nil {
		panic(err)
	}
	s, _ := scanner.New(f)
	tokens, err := s.Scan(f.Name())
	if err != nil {
		panic(err)
	}

	p := parser.New(tokens)
	tree, err := p.Parse(f.Name())
	if err != nil {
		panic(err)
	}

	for _, statement := range tree.Statements {
		content, err := json.MarshalIndent(statement, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(content))
	}
}
