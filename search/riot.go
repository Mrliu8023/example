package main

import (
	"flag"
	"fmt"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var (
	// searcher 是协程安全的
	searcher = riot.Engine{}
)

var file = flag.String("filepath", "groups_dev.json", "groups_dev.json的路径")
var port = flag.String("port", ":18080", "http port")

var useDisk = flag.Bool("use_disk", true, "是否存储")
var storeFolder = flag.String("store_folder", "./riot-index", "store path")
var dictPath = flag.String("dict_path", "./directory.txt", "dict_path")

func main() {
	initSearcher()
	defer searcher.Close()

	//flag.Parse()
	//gl, err := groups_dev.Parse(*file)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, g := range gl.Groups {
	//	searcher.IndexDoc(g.ID, types.DocData{Content: g.Display})
	//}
	//fmt.Println("flush start...")
	//t0 := time.Now()
	//searcher.Flush()
	//fmt.Printf("flush success, spend: %+v\n", time.Since(t0))

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query()["text"]
		if len(text) < 1 {
			w.Write([]byte(`{"err": "wrong text"}`))
			return
		}
		sea := searcher.Search(types.SearchReq{
			Text: text[0],
		})

		w.Write([]byte(fmt.Sprintf("search docs len: %d \n", sea.NumDocs)))

		for _, d := range sea.Docs.(types.ScoredDocs) {
			w.Write([]byte(fmt.Sprintf("id: %s, value: %s\n", d.DocId, d.Content)))
		}
		w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		return
	})

	if err := http.ListenAndServe(*port, nil); err != nil {
		panic(err)
	}
}

func initSearcher() {
	searcher.Init(types.EngineOpts{
		// Using: 1,
		PinYin: true,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.DocIdsIndex,
		},
		NotUseGse:   false,
		UseStore:    *useDisk,
		StoreFolder: *storeFolder,
		GseDict:     *dictPath,
		// GseMode: true,
	})
	os.MkdirAll(*storeFolder, 0777)
}
