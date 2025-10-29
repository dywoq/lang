package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dywoq/lang/pkg/module"
)

func main() {
	f, err := os.Open("./main.dl")
	if err != nil {
		panic(err)
	}
	m := module.MapParser{}
	maps, err := m.Parse(f)
	if err != nil {
		panic(err)
	}
	for key, value := range maps {
		content, _ := json.MarshalIndent(value, "", "  ")
		fmt.Printf("%s - %s\n", key, string(content))
	}
}
