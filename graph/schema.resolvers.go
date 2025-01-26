package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"backoffice/graph/model"
	"context"
	"fmt"
	"crypto/rand" // For generating cryptographically secure random numbers
	"math/big"    // For handling big integers
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
	post := &model.Post{
		ID: fmt.Sprintf("T%d", randNumber),
		Published: false,
		Title: input.Title,
		Text: input.Text,
	}
	r.posts = append(r.posts, post)
	return post, nil
}

// Attach is the resolver for the attach field.
func (r *mutationResolver) Attach(ctx context.Context, files []string) (string, error) {
	panic(fmt.Errorf("not implemented: Attach - attach"))
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.posts, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
