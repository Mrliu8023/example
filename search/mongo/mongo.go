package mongo

import (
	"context"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	_ "go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

type Mongo struct {
	Dbc    *mongo.Client
	DBName string
	sess   mongo.Session
}

type MongoConf struct {
	TimeoutSeconds int
	MaxOpenConns   uint64
	MaxIdleConns   uint64
	URI            string
	RsName         string
	SocketTimeout  int
}

func NewMgo(config MongoConf, timeout time.Duration) (*Mongo, error) {
	connStr, err := connstring.Parse(config.URI)
	if nil != err {
		return nil, err
	}
	//if config.RsName == "" {
	//	return nil, fmt.Errorf("mongodb rsName not set")
	//}
	socketTimeout := time.Second * time.Duration(config.SocketTimeout)
	// do not change this, our transaction plan need it to false.
	// it's related with the transaction number(eg txnNumber) in a transaction session.
	disableWriteRetry := false
	conOpt := options.ClientOptions{
		MaxPoolSize:    &config.MaxOpenConns,
		MinPoolSize:    &config.MaxIdleConns,
		ConnectTimeout: &timeout,
		SocketTimeout:  &socketTimeout,
		ReplicaSet:     &config.RsName,
		RetryWrites:    &disableWriteRetry,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI), &conOpt)
	if nil != err {
		return nil, err
	}

	if err := client.Connect(context.TODO()); nil != err {
		return nil, err
	}

	// TODO: add this check later, this command needs authorize to get version.
	// if err := checkMongodbVersion(connStr.Database, client); err != nil {
	// 	return nil, err
	// }

	// initialize mongodb related metrics
	// initMongoMetric()

	return &Mongo{
		Dbc:    client,
		DBName: connStr.Database,
	}, nil
}