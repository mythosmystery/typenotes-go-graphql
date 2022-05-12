package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	db   *mongo.Database
	Note *mongo.Collection
	User *mongo.Collection
}

func Connect() *DB {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return GetDB(client)
}

func GetDB(client *mongo.Client) *DB {
	return &DB{
		db:   client.Database(os.Getenv("MONGO_DB")),
		Note: client.Database(os.Getenv("MONGO_DB")).Collection("Note"),
		User: client.Database(os.Getenv("MONGO_DB")).Collection("User"),
	}
}
