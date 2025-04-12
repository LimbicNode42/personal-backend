package datasets

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"
	"fmt"
)

// AppendToJSONL is the resolver for the appendToJSONL field.
func (r *mutationResolver) AppendToJSONL(ctx context.Context, fileName string, record map[string]any) (bool, error) {
	panic(fmt.Errorf("not implemented: AppendToJSONL - appendToJSONL"))
}

// ReadJSONL is the resolver for the readJSONL field.
func (r *queryResolver) ReadJSONL(ctx context.Context, fileName string) ([]map[string]any, error) {
	panic(fmt.Errorf("not implemented: ReadJSONL - readJSONL"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
