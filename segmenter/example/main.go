package main

import (
	groups_dev "example/groups-dev"
	"example/search"
	"example/segmenter"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"net/http"
	"os"
	"time"
)

var file = flag.String("filepath", "groups_dev.json", "groups_dev.json的路径")
var port = flag.String("port", ":18080", "http port")

var useDisk = flag.Bool("use_disk", true, "是否存储")
var storeFolder = flag.String("store_folder", "./riot-index", "store path")
var dictPath = flag.String("dict_path", "./dictionary.txt", "dict_path")

var searchMode = flag.Bool("searchMode", true, "searchMode")
var hmm = flag.Bool("hmm", true, "hmm")
var pinyin = flag.Bool("pinyin", true, "pinyin")

func main() {
	flag.Parse()
	fmt.Println("dictPath: ", *dictPath)
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	initSearcher()

	f, _ := os.OpenFile("mem.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	addDocs(*file)

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query()["text"]
		if len(text) < 1 {
			w.Write([]byte(`{"err": "wrong text"}`))
			return
		}
		sea, err := searcher.Search(text[0])
		if err != nil {
			panic(err)
		}

		w.Write([]byte(fmt.Sprintf("search docs len: %d \n", sea.Length)))

		for _, d := range sea.Groups {
			w.Write([]byte(fmt.Sprintf("id: %s, value: %s\n", d.Name(), d.Content())))
		}
		w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		return
	})

	if err := http.ListenAndServe(*port, nil); err != nil {
		panic(err)
	}
}

var searcher search.Searcher

func initSearcher() {
	seg, err := segmenter.NewSegmenter(*searchMode, *hmm, *pinyin, *dictPath)
	if err != nil {
		panic(err)
	}
	searcher = seg
}

func addDocs(filepath string) {
	gl, err := groups_dev.Parse(filepath)
	if err != nil {
		panic(err)
	}

	var docs = make([]search.Document, 0, len(gl.Groups))

	for _, g := range gl.Groups {
		docs = append(docs, g)
	}
	fmt.Println("flush start...")
	t0 := time.Now()
	if err := searcher.AddDocuments(docs); err != nil {
		panic(err)
	}
	fmt.Printf("flush success, spend: %+v\n", time.Since(t0))
}
