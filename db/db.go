package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient holds the MongoDB connection
type MongoClient struct {
	Client *mongo.Client
}

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient(uri string) *MongoClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping MongoDB to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	return &MongoClient{Client: client}
}

// GetCollection returns a MongoDB collection
func (m *MongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}