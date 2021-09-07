package main

import (
	"context"
	"example/search/mongo"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestTractorResource(t *testing.T) {
	conf := mongo.Config{
		Connect: "mongodb://127.0.0.1:27017/dcs_cmdb",
		// RsName: "xbrother",
	}

	db, err := conf.GetMongoClient()
	if err != nil {
		panic(err)
	}

	msClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: "http://127.0.0.1:7700",
	})

	filter := bson.D{{"attributes.location", bson.D{{"$regex", "project_root.*"}}},
			{"deleted", 0},
			{"attributes.ci_type", bson.D{{"$in", []string{"2", "3", "5"}}}}}

	count, err := db.Dbc.Database(db.DBName).Collection("resources").CountDocuments(context.TODO(), filter)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("count: ", count)
	cur, err := db.Dbc.Database(db.DBName).Collection("resources").Find(context.TODO(), filter,
		options.Find().SetProjection(bson.D{{"_id", 1},
			{"attributes.location", 1},
			{"attributes.event_rules", 1},
			{"attributes.ci_type", 1},
			{"attributes.precision", 1},
			{"attributes.value_type", 1},
			{"attributes.mapper", 1},
			{"attributes.compress", 1},
			{"attributes.filter", 1},
			{"attributes.converter", 1},
			{"attributes.spot_type", 1},
			{"attributes.import_type", 1},
			{"version", 1},
			{"attributes.name", 1}}))
	if err != nil {
		panic(err)
	}
	fmt.Println(cur.RemainingBatchLength(), cur.TryNext(context.TODO()))

	var hasDeal = 0
	var rs = make([]map[string]interface{}, 0)

	for cur.Next(context.TODO()) && hasDeal < 200000 {
		 var r map[string]interface{}

		 if err := cur.Decode(&r); err != nil {
		 	panic(err)
		 }

		r["id"] = r["_id"]
		r["revision"] = r["version"]
		delete(r, "_id")
		delete(r, "version")

		if len(rs) < 100 {
			rs = append(rs, r)
			hasDeal++
			continue
		}
		_, err = msClient.Index("resources").AddDocumentsWithPrimaryKey(rs, "id")
		if err != nil {
			panic(err)
		}

		rs = make([]map[string]interface{}, 0)

		hasDeal++
	}

	cur.Close(context.TODO())
}
