package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
)

type A struct {
	a bool
	b int64
	c bool
}

type B struct {
	A bool
	B bool
	C int64
}

// input 1.2.17
// output 1 1.2 1.2.17
func main() {
	b := B{
		A: true,
		B: true,
		C: 134,
	}

	bs, _ := json.Marshal(b)
	dj := datatypes.JSON(bs)

	var c B
	json.Unmarshal([]byte(dj), &c)
	fmt.Println(c)
}
