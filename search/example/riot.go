package main

import (
	"fmt"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"os"
)

var (
	// searcher 是协程安全的
	searcher = riot.Engine{}
)

func initEngine() {
	var path = "./riot-index"

	searcher.Init(types.EngineOpts{
		Using:   4,
		GseMode: true,
		PinYin:  true,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.DocIdsIndex,
		},
		UseStore: false,
		// StoreFolder: path,
		// NotUseGse: true,
		GseDict: "dictionary.txt",
		// StopTokenFile:           "../../riot/data/dict/stop_tokens.txt",
	})
	defer searcher.Close()
	os.MkdirAll(path, 0777)

	//text := "在路上, in the way"
	//index1 := types.DocData{Content: text}
	//index2 := types.DocData{Content: text}
	//index3 := types.DocData{Content: "In the way."}
	//index4 := types.DocData{Content: "温湿度TH03."}
	//index5 := types.DocData{Content: "基础设施/综合布线/其他/维谛/维谛/公共模块/温湿度", Tokens: []types.TokenData{types.TokenData{Text: "温湿度", Locations: []int{
	//	strings.LastIndex("基础设施/综合布线/其他/维谛/维谛/公共模块/温湿度", "/") + 1,
	//}}}}

	fmt.Println(searcher.Segment("基础设施/综合布线/其他/维谛/维谛/公共模块/温湿度Th03"))
	fmt.Println(searcher.PinYin("基础设施/综合布线/其他/维谛/维谛/公共模块/温湿度"))

	//searcher.Index("10", index1)
	//searcher.Index("11", index2)
	//searcher.Index("12", index3)
	//searcher.Index("13", index4)
	//searcher.Index("14", index5)
	////
	//// 等待索引刷新完毕
	//searcher.Flush()
}

func main() {
	initEngine()

	//sea := searcher.Search(types.SearchReq{
	//	Text: "wsd",
	//})
	//
	//fmt.Println("search response: ", sea, "; docs = ", sea.Docs)
}
