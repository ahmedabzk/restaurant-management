package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dbinstance() *mongo.Client{
	// load the .env file
	if err := godotenv.Load(".env"); err != nil{
		fmt.Println("no .env file found")
	}
	// get the uri from the .env file
	uri := os.Getenv("MONGOURL")
	if uri == ""{
		log.Fatal("you must set mongo url")
	}
	// create new client
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil{
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("connected to mongodb")

	return client
}

var Client *mongo.Client = Dbinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	return collection
}