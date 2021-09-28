package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func TestQuery(t *testing.T) {
	c := redisearch.NewClient("192.168.30.58:6379", "groups")
	docs, total, err := c.Search(redisearch.NewQuery("th03"))
	if err != nil {
		log.Fatal(err)
	}

	for _, doc := range docs {
		fmt.Println(doc.Id, doc.Properties)
	}
	fmt.Println("total: ", total)
}
