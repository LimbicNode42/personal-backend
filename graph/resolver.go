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
	SMBClient *db.SMBClient
	posts []*model.Post
}

func NewResolver(mongoURI string, smbConfig *db.SMBConfig) *Resolver {
	mongoClient := db.NewMongoClient(mongoURI)

	smbClient, err := db.SMBConnect(smbConfig)
	if err != nil {
		log.Fatalf("SMB connection failed: %v", err)
	}

	var posts []*model.Post

	return &Resolver{
		MongoClient: mongoClient,
		SMBClient: smbClient,
		posts: posts,
	}
}

func (r *Resolver) Close() {
	if r.MongoClient != nil {
		r.MongoClient.Close()
	}
	if r.SMBClient != nil {
		r.SMBClient.SMBClose()
	}
}