package db

import (
    "context"
    "log"
	// "os"
    "time"

	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

	// "backoffice/auth"
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

	log.Println("Connected to MongoDB")

	return &MongoClient{Client: client}
}

// GetCollection returns a MongoDB collection
func (m *MongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.Client.Database(database).Collection(collection)
}

func CreateMongoUri() string {
	//TODO: add auth to mongodb
	// see https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Connect
	// client := auth.InfisicalLogin()
	// log.Println("Project ID: %v", os.Getenv("INF_DEV_PROJECT_ID"))
	// secrets := auth.InfisicalGetSecrets(client,os.Getenv("INF_DEV_PROJECT_ID"),"prod","/mongo")

	uri := "mongodb://192.168.0.111:27017"

	return uri
}

// Function to get next incrementing ID
func GetNextCollectionIndex(collection *mongo.Collection, counterName string) (int, error) {
	filter := bson.M{"_id": counterName}
	update := bson.M{"$inc": bson.M{"index": 1}}

	// Find and update the counter atomically
	var result struct {
		Index int `bson:"index"`
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)

	if err != nil {
		return 0, err
	}
	return result.Index, nil
}