package main

import (
	"context"
	"kenja2tools"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type copy struct {
	ctx           context.Context
	mongoClient   *mongo.Client
	opnsrchClient *opensearchapi.Client

	database   *mongo.Database
	collection *mongo.Collection
}

func (run *copy) run() error {
	f := bson.M{}
	op := options.Find()
	stream, err := run.collection.Find(run.ctx, f, op)
	if err != nil {
		return err
	}
	defer stream.Close(run.ctx)

	for stream.Next(run.ctx) {
		doc := bson.M{}
		if err := stream.Decode(&doc); err != nil {
			return err
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}

	log.Println("done")
	return nil
}

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Panicln("failed to load .env file:", err)
	}
	ctx := context.Background()

	opnsrchClient, err := kenja2tools.OpensearchApiClient()
	if err != nil {
		log.Panicln("failed to connect to opensearch:", err)
	}

	mongoClient, err := kenja2tools.MongoClient()
	if err != nil {
		log.Panicln("failed to connect to mongo:", err)
	}
	defer mongoClient.Disconnect(ctx)

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
		ctx,
		mongoClient,
		opnsrchClient,
		database,
		collection,
	}
	if err := c.run(); err != nil {
		log.Fatalln(err)
	}
}
