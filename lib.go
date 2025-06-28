package kenja2tools

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"kenja2tools/J"
	"log"
	"net/http"
	"os"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoConfig() *options.ClientOptions {
	uri := os.Getenv("MONGO_URI")
	if len(uri) == 0 {
		log.Panicln("env for mongo uri is not set")
	}

	return options.Client().ApplyURI(uri)
}

func NewMongoClient() (*mongo.Client, error) {
	cfg := NewMongoConfig()
	return mongo.Connect(cfg)
}

func NewOpensearchConfig() opensearch.Config {
	uri := os.Getenv("OPENSEARCH_URI")
	if len(uri) == 0 {
		log.Panicln("env for opensearch uri is not set")
	}
	name := os.Getenv("OPENSEARCH_NAME")
	if len(name) == 0 {
		log.Panicln("env for opensearch name is not set")
	}
	pw := os.Getenv("OPENSEARCH_PW")
	if len(pw) == 0 {
		log.Panicln("env for opensearch pw is not set")
	}

	return opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{uri},
		Username:  name,
		Password:  pw,
	}
}

func NewOpensearchClient() (*opensearch.Client, error) {
	cfg := NewOpensearchConfig()
	return opensearch.NewClient(cfg)
}

func NewOpensearchApiClient() (*opensearchapi.Client, error) {
	cfg := opensearchapi.Config{Client: NewOpensearchConfig()}
	return opensearchapi.NewClient(cfg)
}

func NewBodyFromDocument(doc bson.M) (string, io.Reader, error) {
	id, ok := doc["_id"]
	if !ok {
		return "", nil, errors.New("could not find _id in document")
	}
	oid, ok := id.(bson.ObjectID)
	if !ok {
		return "", nil, errors.New("_id is not an ObjectID")
	}
	delete(doc, "_id")

	b, err := json.Marshal(doc)
	if err != nil {
		return "", nil, err
	}

	return oid.Hex(), bytes.NewReader(b), nil
}

func NewBulkIndexReqBody(index string, doc bson.M) (io.Reader, error) {
	id, body, err := NewBodyFromDocument(doc)
	if err != nil {
		return nil, err
	}

	req := J.Json{
		"create": J.Json{
			"_index": index,
			"_id":    id,
		},
	}
	reqBody, err := req.Reader()
	if err != nil {
		return nil, err
	}

	return io.MultiReader(reqBody, body), nil
}

func NewIndexReqsFromDocuments(index string, docs []bson.M) (opensearchapi.BulkReq, error) {
	bulk := []io.Reader{}

	for _, doc := range docs {
		body, err := NewBulkIndexReqBody(index, doc)
		if err != nil {
			log.Println("skipping...", err)
			continue
		}

		bulk = append(bulk, body)
	}

	return opensearchapi.BulkReq{
		Index: index,
		Body:  io.MultiReader(bulk...),
	}, nil
}
