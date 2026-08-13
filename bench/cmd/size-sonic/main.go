package main

import (
	"fmt"

	"github.com/bytedance/sonic"
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

	b, err := sonic.ConfigFastest.Marshal(value)
	fmt.Println(string(b), err)
}
