package main

import (
	"example/groups-dev"
	"fmt"
	"github.com/go-ego/riot/types"
	"testing"
	"time"
)

func TestRiot(t *testing.T) {
	initSearcher()
	gl, err := groups_dev.Parse("groups_dev.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, g := range gl.Groups {
		searcher.IndexDoc(g.ID, types.DocData{Content: g.String()})
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
