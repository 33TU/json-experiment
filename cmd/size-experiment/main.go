package main

import (
	"fmt"

	jsonexperiment "github.com/33TU/json-experiment"
)

func main() {
	var value = struct {
		Name    string         `json:"name"`
		Numbers []int          `json:"numbers"`
		Meta    map[string]any `json:"meta"`
	}{
		Name:    "benchmark",
		Numbers: []int{1, 2, 3},
		Meta:    map[string]any{"ok": true},
	}

	b, err := jsonexperiment.Marshal(value)
	fmt.Println(string(b), err)
}
