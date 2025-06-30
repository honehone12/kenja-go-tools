package main

import (
	"context"
	"errors"
	"flag"
	"kenja2tools"
	"kenja2tools/documents"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type copy struct {
	ctx       context.Context
	batchSize int

	mongoClient   *mongo.Client
	opnsrchClient *opensearchapi.Client

	database   *mongo.Database
	collection *mongo.Collection

	index string
}

func (run *copy) bulkRequest(batch []documents.FlatDocument) error {
	bulk, err := kenja2tools.NewIndexReqsFromDocuments(run.index, batch)
	if err != nil {
		return err
	}

	res, err := run.opnsrchClient.Bulk(run.ctx, bulk)
	if err != nil {
		return err
	}

	if res.Errors {
		for _, item := range res.Items {
			for op, info := range item {
				if info.Error != nil {
					log.Printf("[ERROR] %s: %#v\n", op, info.Error)
				}
			}
		}
		return errors.New("failed to copy items, see above for details")
	} else {
		log.Printf("copied %d items\n", len(res.Items))
		return nil
	}
}

func (run *copy) run() error {
	f := bson.M{}
	op := options.Find()
	stream, err := run.collection.Find(run.ctx, f, op)
	if err != nil {
		return err
	}
	defer stream.Close(run.ctx)

	batch := make([]documents.FlatDocument, 0, run.batchSize)

	for stream.Next(run.ctx) {
		doc := documents.FlatDocument{}
		if err := stream.Decode(&doc); err != nil {
			return err
		}

		batch = append(batch, doc)
		if len(batch) >= run.batchSize {
			if err := run.bulkRequest(batch); err != nil {
				return err
			}
			batch = make([]documents.FlatDocument, 0, run.batchSize)
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}

	if len(batch) > 0 {
		if err := run.bulkRequest(batch); err != nil {
			return err
		}
	}

	log.Println("done")
	return nil
}

func args() int {
	batchSize := flag.Int("batch-size", 100, "number of documents to copy in one batch")
	flag.Parse()

	return *batchSize
}

func main() {
	batchSize := args()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Panicln("failed to load .env file:", err)
	}
	ctx := context.Background()

	opnsrchClient, err := kenja2tools.NewOpensearchApiClient()
	if err != nil {
		log.Panicln("failed to connect to opensearch:", err)
	}

	mongoClient, err := kenja2tools.NewMongoClient()
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

	index := os.Getenv("OPENSEARCH_IDX")
	if len(index) == 0 {
		log.Panicln("env for opensearch index is not set")
	}

	c := copy{
		ctx,
		batchSize,
		mongoClient,
		opnsrchClient,
		database,
		collection,
		index,
	}
	if err := c.run(); err != nil {
		log.Fatalln(err)
	}
}
