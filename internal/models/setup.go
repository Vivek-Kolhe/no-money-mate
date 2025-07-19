package models

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectDB() *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGOURI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connection err: ", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping err: ", err)
	}

	log.Println("MongoDB connected")

	dbName := os.Getenv("DB_NAME")
	return &Database{
		Client:   client,
		Database: client.Database(dbName),
	}
}

func (db *Database) Disconnect(ctx context.Context) {
	if err := db.Client.Disconnect(ctx); err != nil {
		log.Panic("MongoDB disconnect err: ", err)
		return
	}
	log.Println("MongoDB disconnected successfully")
}

func (db *Database) GetCollection(collectionName string) *mongo.Collection {
	return db.Database.Collection(collectionName)
}
