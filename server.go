package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	blog_api "backoffice/graph/blog"
	datasets_api "backoffice/graph/datasets"
	"backoffice/auth"
	"backoffice/db"
)

const defaultPort = "8080"

var KeycloakURL = "https://192.168.0.109:8443"
var KeycloakRealm = "shadow"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize JWKS for Keycloak authentication
	if err := auth.InitJWKS(); err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	cdnSecrets := auth.InfisicalGetSecrets("/omv")

	mongoURI := db.CreateMongoUri()
	smbConfig, err := db.SMBConfigure(cdnSecrets)
	if err != nil {
		panic(err)
	}
	blog_resolver := blog_api.NewResolver(mongoURI, smbConfig)
	defer blog_resolver.Close()
	blog := handler.New(blog_api.NewExecutableSchema(blog_api.Config{Resolvers: blog_resolver}))
	blog.AddTransport(transport.Options{})
	blog.AddTransport(transport.GET{})
	blog.AddTransport(transport.POST{})
	blog.AddTransport(transport.MultipartForm{})
	blog.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	blog.Use(extension.Introspection{})
	blog.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	datasets_resolver := datasets_api.NewResolver()
	defer datasets_resolver.Close()
	datasets := handler.New(datasets_api.NewExecutableSchema(datasets_api.Config{Resolvers: datasets_resolver}))
	datasets.AddTransport(transport.Options{})
	datasets.AddTransport(transport.GET{})
	datasets.AddTransport(transport.POST{})
	datasets.AddTransport(transport.MultipartForm{})
	datasets.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	datasets.Use(extension.Introspection{})
	datasets.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Set up routerication middleware
	router := chi.NewRouter()

	log.Println("Applying Logging Middleware")
	router.Use(middleware.Logger)
	
	// Enable CORS
	log.Println("Initializing CORS middleware")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001", "https://backoffice.wheeler-network.com"}, // Allow React app
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	router.Use(corsHandler.Handler)

	log.Println("Applying Auth Middleware")
	router.Use(auth.AuthMiddleware) // Apply auth to all routes except `/playground`

	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/blog", blog)
	router.Handle("/datasets", datasets)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		// TODO: handle error better
		panic(err)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
