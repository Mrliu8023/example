package groups_dev

import (
	"fmt"
	"testing"
)

func Test_parse(t *testing.T) {
	gl, err := Parse("groups_dev.json")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(gl.Length)

	for i, g := range gl.Groups {
		if i > 10 {
			return
		}
		fmt.Println(g.String())
	}
}
