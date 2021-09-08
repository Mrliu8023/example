package groups_dev

import (
	"bufio"
	"fmt"
	"os"
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
		fmt.Println(g.Content())
	}
}

func TestBuildDict(t *testing.T) {
	gl, err := Parse("D:\\go\\src\\example\\groups_dev.json")
	if err != nil {
		t.Fatal(err)
	}

	fp, err := os.OpenFile("D:\\go\\src\\example\\groups-dev\\dict_groups.txt", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)
	for _, g := range gl.Groups {
		w.WriteString(fmt.Sprintf("%s\n", g.Value))
	}
	w.Flush()
}
