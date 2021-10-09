package main

import (
	groups_dev "example/groups-dev"
	"example/search"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	c := getMSClient()

	gl, err := groups_dev.Parse("D:\\go\\src\\example\\groups_dev.json")
	if err != nil {
		panic(err)
	}

	indexConfig := &meilisearch.IndexConfig{
		Uid:        "groups-dev",
		PrimaryKey: "id",
	}

	c.DeleteIndex("groups-dev")

	_, err = c.CreateIndex(indexConfig)
	if err != nil {
		if !strings.Contains(err.Error(), "index_already_exists") {
			panic(err)
		}
	}

	docs := []map[string]interface{}{}

	// a document primary key can be of type integer or string only composed of alphanumeric characters, hyphens (-) and underscores (_).
	for _, g := range gl.Groups {
		doc := map[string]interface{}{
			"id":     strings.ReplaceAll(g.ID, ".", "_"),
			"value":  g.Content(),
			"pinyin": search.PinYin(strings.ReplaceAll(g.Display, "/", ",")),
		}
		docs = append(docs, doc)
	}

	id, err := c.Index("groups-dev").AddDocuments(docs)
	if err != nil {
		panic(err)
	}
	t0 := time.Now()
	for {
		resp, err := c.Index("groups-dev").GetUpdateStatus(id.UpdateID)
		if err != nil {
			panic(err)
		}
		if resp.Status == meilisearch.UpdateStatusProcessed {
			fmt.Println("update ok!, spend: ", time.Since(t0))
			break
		} else if resp.Status == meilisearch.UpdateStatusFailed {
			panic(fmt.Errorf("update failed: %s", resp.Error))
		}
		time.Sleep(3 * time.Second)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query()["text"]
		if len(text) < 1 {
			w.Write([]byte(`{"err": "wrong text"}`))
			return
		}
		fmt.Println("text: ", text[0])
		searchRes, err := c.Index("groups-dev").Search(text[0],
			&meilisearch.SearchRequest{
				Limit: 10,
			})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(searchRes.Hits)
		w.Write([]byte(fmt.Sprintf("search docs len: %d \n", len(searchRes.Hits))))

		for _, d := range searchRes.Hits {
			w.Write([]byte(fmt.Sprintf("value: %+v\n", d)))
		}
		w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		return
	})

	http.ListenAndServe(":8081", nil)

}
