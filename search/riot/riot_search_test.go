package main

import (
	"encoding/json"
	groups_dev "example/groups-dev"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-ego/riot/types"
)

func TestRiot(t *testing.T) {
	initSearcher()
	gl, err := groups_dev.Parse("groups_dev.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, g := range gl.Groups {
		searcher.IndexDoc(g.ID, types.DocData{Content: g.Content()})
	}
	fmt.Println("flush start...")
	t0 := time.Now()
	searcher.Flush()
	fmt.Printf("flush success, spend: %+v\n", time.Since(t0))

	sea := searcher.Search(types.SearchReq{
		Text: "wsd",
	})

	fmt.Println("search response: ", sea)

	for _, d := range sea.Docs.(types.ScoredDocs) {
		fmt.Println("id: ", d.DocId, "; value: ", d.Content)
	}

}

func TestAddDoc(t *testing.T) {
	filepath := "D:\\go\\src\\example\\groups_dev.json"
	gl, err := groups_dev.Parse(filepath)
	if err != nil {
		panic(err)
	}

	for _, g := range gl.Groups {
		bs, _ := json.Marshal(g)
		resp, err := http.Post("http://192.168.30.58:18080/group", "application/json", strings.NewReader(string(bs)))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not ok")
		}
		fmt.Printf("success: %+v\n", g)
	}
}
