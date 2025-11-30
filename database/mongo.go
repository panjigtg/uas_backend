package database

import (
	"context"
	"log"
	"time"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func MongoConnections() *mongo.Database {	

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("Set Url Mongo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Can't Connect:", err)
	}


	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Can't Connect:", err)
	}

	log.Println("mongo connect")

	return client.Database(os.Getenv("MONGO_DB"))
}