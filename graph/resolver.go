package graph

import (
	"log"

	"backoffice/graph/model"
	"backoffice/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	MongoClient *db.MongoClient

	posts []*model.Post
}

func NewResolver() *Resolver {
	mongoURI := "mongodb://192.168.0.110:27017" // Your MongoDB URI
	mongoClient := db.NewMongoClient(mongoURI)

	log.Println("Connected to MongoDB")

	var posts []*model.Post

	return &Resolver{
		MongoClient: mongoClient,
		posts: posts,
	}
}