package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/dywoq/lang/pkg/parser"
	"github.com/dywoq/lang/pkg/scanner"
)

func main() {
	file := flag.String("file", "main.dl", "A file to interpret")
	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	s, err := scanner.NewDebug(f, os.Stdout)
	if err != nil {
		panic(err)
	}

	tokens, err := s.Scan()
	if err != nil {
		panic(err)
	}
	for _, token := range tokens {
		fmt.Printf("%s -- %s\n", token.Literal, token.Kind)
	}

	p := parser.NewDebug(tokens, os.Stdout)

	tree, err := p.Parse()
	if err != nil {
		panic(err)
	}

	content, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("string(content): %v\n", string(content))
}
