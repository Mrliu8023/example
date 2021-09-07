package mongo

import (
	"fmt"
	"testing"
)

func TestNewMgo(t *testing.T) {
	conf := Config{
		Connect: "mongodb://127.0.0.1:27017/cmdb",
		// RsName: "xbrother",
	}

	db, err := conf.GetMongoClient()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", db)
}
