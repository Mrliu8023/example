package main

import (
	"fmt"
	"strings"
)

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
	a := "1"
	fmt.Println(a[:strings.LastIndex(a, ".")])
}
