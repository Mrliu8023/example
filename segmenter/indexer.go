package segmenter

import (
	"sync"
)

type indexerAddDocReq struct {
	doc *DocIndex
}

type Indexer struct {
	sync.RWMutex
	table map[string][]string

	docCache map[string]*Document
}

// DocIndex document's index
type DocIndex struct {
	// DocId 文本的 DocId
	DocId string
	//
	//// TokenLen 文本的关键词长
	//TokenLen float32

	// Keywords 加入的索引键
	Keywords []string
}

// KeywordIndex 反向索引项，这实际上标注了一个（搜索键，文档）对。
type KeywordIndex struct {
	// Text 搜索键的 UTF-8 文本
	Text string
	//
	//// Frequency 搜索键词频
	//Frequency float32
	//
	//// Starts 搜索键在文档中的起始字节位置，按照升序排列
	//Starts []int
}

func (indexer *Indexer) AddIndex(req *indexerAddDocReq) {
	indexer.RWMutex.Lock()
	defer indexer.RWMutex.Unlock()

	for _, Keyword := range req.doc.Keywords {
		indexer.table[Keyword] = append(indexer.table[Keyword], req.doc.DocId)
	}
	// TODO 相同的ID 怎么处理? 将id排序（二叉搜索树？）
}

func (indexer *Indexer) Lookup(keyword string) []string {
	indexer.RLock()
	defer indexer.RUnlock()
	return indexer.table[keyword]
}
