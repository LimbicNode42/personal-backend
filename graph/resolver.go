package graph

import (
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
	mongoURI := db.CreateMongoUri()
	mongoClient := db.NewMongoClient(mongoURI)

	var posts []*model.Post

	return &Resolver{
		MongoClient: mongoClient,
		posts: posts,
	}
}