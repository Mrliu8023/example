package main

import (
	"encoding/json"
	groups_dev "example/groups-dev"
	"flag"
	"fmt"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/pkg/profile"
	"io/ioutil"
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
var dictPath = flag.String("dict_path", "./dict.txt", "dict_path")
var storeEngine = flag.String("storeEngine", "bg", "storeEngine")

func main() {
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	initSearcher()
	defer searcher.Close()

	f, _ := os.OpenFile("mem.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	flag.Parse()
	//addDocs(*file)

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

	http.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		var g = &groups_dev.Group{}
		if err := json.Unmarshal(bs, g); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		searcher.IndexDoc(g.ID, types.DocData{Content: g.Display}, false)
		searcher.Flush()

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
		StoreFolder: *storeFolder,
		GseDict:     *dictPath,
		GseMode:     true,
		Hmm:         true,
		StoreEngine: *storeEngine,
		StoreShards: 2,
	})
	os.MkdirAll(*storeFolder, 0777)
}

func addDocs(filepath string) {

}
