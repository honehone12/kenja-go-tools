package main

import (
	"kenja2tools"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type copy struct {
	mongoClient   *mongo.Client
	opnsrchClient *opensearch.Client

	database   *mongo.Database
	collection *mongo.Collection
}

func (op *copy) run() error {

	log.Println("done")
	return nil
}

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Panicln("failed to load .env file:", err)
	}

	mongoClient, err := kenja2tools.MongoClient()
	if err != nil {
		log.Panicln("failed to connect to mongo:", err)
	}

	opnsrchClient, err := kenja2tools.OpensearchClient()
	if err != nil {
		log.Panicln("failed to connect to opensearch:", err)
	}

	db := os.Getenv("MONGO_DB_NAME")
	if len(db) == 0 {
		log.Panicln("env for mongo db name is not set")
	}
	database := mongoClient.Database(db)

	cl := os.Getenv("MONGO_CL_NAME")
	if len(cl) == 0 {
		log.Panicln("env for mongo cl name is not set")
	}
	collection := database.Collection(cl)

	c := copy{
		mongoClient,
		opnsrchClient,
		database,
		collection,
	}
	if err := c.run(); err != nil {
		log.Fatalln(err)
	}
}
