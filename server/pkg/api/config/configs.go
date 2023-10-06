package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func DBInstance() {
	mongodbAtlasUri := os.Getenv("MONGODB_ATLAS_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbAtlasUri))
	if err != nil {
		log.Fatal("🔥 failed to connect to the Database!\n", err.Error())
	}

	log.Println("🚀 Connected Successfully to the Database")

	Client = newClient
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dBName := os.Getenv("MONGO_DATABASE_NAME")

	return client.Database(dBName).Collection(collectionName)
}
