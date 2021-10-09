package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	groups_dev "example/groups-dev"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func main() {
	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClient("192.168.30.58:6380", "groups")

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextField("value")).
		AddField(redisearch.NewTextField("pinyin"))

	// Drop an existing index. If the index does not exist an error is returned
	c.Drop()

	// Create the index with the given schema
	if err := c.CreateIndexWithIndexDefinition(sc, redisearch.NewIndexDefinition().SetLanguage("Chinese")); err != nil {
		log.Fatal(err)
	}

	// if err := c.CreateIndex(sc); err != nil {
	// 	log.Fatal(err)
	// }

	gl, err := groups_dev.Parse("D:\\go\\src\\example\\groups_dev.json")
	if err != nil {
		log.Fatal(err)
	}

	var docs = make([]redisearch.Document, 0, gl.Length)
	fmt.Println("gl len: ", gl.Length)
	for _, g := range gl.Groups {
		// Create a document with an id and given score
		doc := redisearch.NewDocument(g.ID, 1.0)
		doc.Set("name", g.ID).
			Set("value", g.Display).
			Set("pinyin", PinYin(strings.ReplaceAll(g.Display, "/", ",")))
		docs = append(docs, doc)
	}

	log.Println("will add docs num: ", len(docs))
	t0 := time.Now()
	// Index the document. The API accepts multiple documents at a time
	if err := c.Index(docs...); err != nil {
		log.Fatal(err)
	}
	log.Println("will add docs succeed: ", time.Since(t0))

	// Searching with limit and sorting
	docs, total, err := c.Search(redisearch.NewQuery("gongji"))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(docs[0].Id, docs[0].Properties, total, err)
	for _, doc := range docs {
		fmt.Println(doc.Id, doc.Properties)
	}
	fmt.Println("total: ", total)
	// Output: doc1 Hello world 1 <nil>
}
