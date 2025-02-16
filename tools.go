//go:build tools

package tools

import (
	_ "github.com/go-resty/resty/v2"
	_ "github.com/go-chi/chi/v5"
	_ "github.com/rs/cors"
	
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"

	_ "go.mongodb.org/mongo-driver/mongo"

	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/MicahParks/keyfunc/v3"
	_ "github.com/MicahParks/jwkset"

	_ "github.com/infisical/go-sdk"
)
