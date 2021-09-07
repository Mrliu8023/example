package main

import (
	"fmt"
	"github.com/meilisearch/meilisearch-go"
)

func main1() {
	client := getMSClient()

	ds := []map[string]interface{}{
		{"id": "5_4_0_108_1_1_0", "name": "温度2222", "test": 1},
	}
	updateID, err := client.Index("resources").UpdateDocuments(ds)
	if err != nil {
		panic(err)
	}

	resp, err := client.Index("resources").GetUpdateStatus(updateID.UpdateID)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func getMSClient() *meilisearch.Client {

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: "http://127.0.0.1:7700",
	})

	return client
}
