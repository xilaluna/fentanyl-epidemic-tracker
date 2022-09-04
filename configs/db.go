package configs

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func ConnectDB() *mongo.Client {
	if client != nil {
		return client
	}
	
	client, err := mongo.NewClient(options.Client().ApplyURI(LoadMongoEnv()))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return client
}

func CloseDB() {
	if client == nil {
		return
	}
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func GetClient() *mongo.Client {
	if client != nil {
		return client
	}
	return ConnectDB()
}


func DatabaseCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("fentanyl-epidemic-data").Collection("articles")
	return collection
}