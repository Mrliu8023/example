package main

import "fmt"

type A struct {
	a bool
	b int64
	c bool
}

type B struct {
	a bool
	b bool
	c int64
}

// input 1.2.17
// output 1 1.2 1.2.17
func main() {
	a := make(map[string][]string)

	a["1"] = append(a["1"], "1", "2")
	fmt.Println(a["1"])
}
