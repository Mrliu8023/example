package segmenter

import (
	"context"
	"crypto/md5"
	groups_dev "example/groups-dev"
	"example/search"
	"fmt"
	"github.com/go-ego/gpy"
	"github.com/go-ego/gse"
	"github.com/go-ego/riot/types"
	"hash"
	"runtime"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Segmenter struct {
	seg gse.Segmenter

	searchMode, hmm, pinyin bool

	indexer *Indexer

	sync.RWMutex
	docCache map[string]search.Document

	rc *redis.Client

	hash hash.Hash
	//// 建立索引器使用的通信通道
	NumGseThreads int
	segmenterChan chan search.Document
	//indexerAddDocChans []chan indexerAddDocReq
}

type segmenterReq struct {
	docId string
	data  types.DocData
}

type Document struct {
	id      string
	content string
}

func (doc *Document) ID() string {
	return doc.id
}

func (doc *Document) Name() string {
	return doc.id
}

func (doc *Document) Content() string {
	return doc.content
}

func NewSegmenter(mode, hmm, pinyin bool, dictPath ...string) (*Segmenter, error) {
	s, err := gse.New(dictPath...)
	if err != nil {
		return nil, err
	}

	rdb := initRedisClient()

	seg := &Segmenter{
		seg:        s,
		searchMode: mode,
		hmm:        hmm,
		pinyin:     pinyin,

		hash: md5.New(),

		indexer:  &Indexer{table: map[string][]string{}},
		docCache: map[string]search.Document{},

		rc: rdb,
	}

	if seg.NumGseThreads == 0 {
		seg.NumGseThreads = runtime.NumCPU()
	}

	seg.segmenterChan = make(
		chan search.Document, seg.NumGseThreads)

	fmt.Println("segmentWorker numbers: ", seg.NumGseThreads)

	for i := 0; i < seg.NumGseThreads; i++ {
		go seg.segmentWorker()
	}

	return seg, nil
}

func (sm *Segmenter) Segment(text string) []string {
	var (
		str      string
		pyStr    string
		strArr   []string
		pyArr    []string
		splitStr string
	)

	//
	splitHans := strings.Split(text, "")
	for i := 0; i < len(splitHans); i++ {
		if splitHans[i] != "" {
			strArr = append(strArr, splitHans[i])
			splitStr += splitHans[i]
			strArr = append(strArr, splitStr)
		}
	}

	// Segment 分词

	sehans := sm.segment(text)
	for h := 0; h < len(sehans); h++ {
		strArr = append(strArr, sehans[h])
	}

	//
	// py := pinyin.LazyConvert(sehans[h], nil)
	pyMap := make(map[string]struct{})
	// fmt.Println(strArr)
	py := gpy.LazyConvert(text, nil)

	// fmt.Println("py...", py)
	for i := 0; i < len(py); i++ {
		// log.Println("py[i]...", py[i])
		pyStr += py[i]

		pyMap[pyStr] = struct{}{}
		pyArr = append(pyArr, pyStr)

		if len(py[i]) > 0 {
			str += py[i][0:1]

			pyMap[pyStr] = struct{}{}
			pyArr = append(pyArr, str)

		}
	}

	for _, han := range strArr {
		str = ""
		py = gpy.LazyConvert(han, nil)
		// fmt.Println("py: ", py)
		for i := 0; i < len(py); i++ {
			if _, ok := pyMap[py[i]]; !ok {
				pyMap[py[i]] = struct{}{}
				pyArr = append(pyArr, py[i])
			}
			if len(py[i]) > 0 {
				str += py[i][0:1]

				if _, ok := pyMap[str]; !ok {
					pyMap[py[i]] = struct{}{}
					pyArr = append(pyArr, str)
				}

			}
		}
		// fmt.Println("pyArr: ", pyArr)
	}
	strArr = append(strArr, pyArr...)

	return strArr
}

func (sm *Segmenter) segment(text string) []string {
	var segments []string

	if sm.searchMode {
		segments = sm.seg.CutSearch(text, sm.hmm)
	} else {
		segments = sm.seg.Cut(text, sm.hmm)
	}
	return segments
}

func (sm *Segmenter) Search(keyword string) (*groups_dev.GroupList, error) {
	key := sm.hash.Sum([]byte(keyword))
	//ids := sm.indexer.Lookup(keyword)
	resp := sm.rc.SMembers(context.Background(), string(key))

	ids, err := resp.Result()
	if err != nil {
		return nil, err
	}

	gl := &groups_dev.GroupList{Groups: make([]*groups_dev.Group, 0, len(ids)), Length: len(ids)}
	for _, id := range ids {
		var g = &groups_dev.Group{}
		g.ID = id
		sm.RLock()
		g.Display = sm.docCache[id].Content()
		sm.RUnlock()
		g.Value = strings.Split(g.Display, "/")[len(strings.Split(g.Display, "/"))-1]
		gl.Groups = append(gl.Groups, g)
	}

	return gl, nil
}

func (sm *Segmenter) AddDocuments(docs []search.Document) error {
	for _, doc := range docs {
		sm.segmenterChan <- doc
	}

	return nil
}

func (sm *Segmenter) segmentWorker() {
	for {
		doc := <-sm.segmenterChan
		segments := sm.Segment(doc.Content())
		for _, segment := range segments {
			key := sm.hash.Sum([]byte(segment))
			if err := sm.rc.SAdd(context.Background(), fmt.Sprintf("%x", key), doc.Name()).Err(); err != nil {
				panic(err)
			}
		}

		sm.Lock()
		sm.docCache[doc.Name()] = doc
		sm.Unlock()
		fmt.Printf("deal doc: %s, value: %s success\n", doc.Name(), doc.Content())
	}
}

var ctx = context.Background()

func initRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.30.58:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 10,
	})
	return rdb

	//err := rdb.Set(ctx, "key", "value", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//val, err := rdb.Get(ctx, "key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)
	//
	//val2, err := rdb.Get(ctx, "key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
	//// Output: key value
	//// key2 does not exist
}
