package kenja2tools

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/opensearch-project/opensearch-go/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func MongoConfig() *options.ClientOptions {
	uri := os.Getenv("MONGO_URI")
	if len(uri) == 0 {
		log.Panicln("env for mongo uri is not set")
	}

	return options.Client().ApplyURI(uri)
}

func MongoClient() (*mongo.Client, error) {
	cfg := MongoConfig()
	return mongo.Connect(cfg)
}

func OpensearchConfig() opensearch.Config {
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

func OpensearchClient() (*opensearch.Client, error) {
	cfg := OpensearchConfig()
	return opensearch.NewClient(cfg)
}
