package auth

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/MicahParks/jwkset"
)

var KeycloakURL = "https://192.168.0.109:8443"
var KeycloakRealm = "shadow"

var keycloakJWKS keyfunc.Keyfunc

// Initialize JWKS from Keycloak
func InitJWKS() (error) {
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", KeycloakURL, KeycloakRealm)
	log.Println("Fetching JWKS from:", jwksURL)

	// Create an HTTP client that ignores TLS verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ⚠️ Not safe for production
	}
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	storage, err := jwkset.NewStorageFromHTTP(jwksURL, jwkset.HTTPClientStorageOptions{
		Client:          client,
		RefreshInterval: time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to create JWKS storage: %v", err)
		return err
	}

	ctx := context.Background()
	keycloakJWKS, err = keyfunc.New(keyfunc.Options{
		Ctx:     ctx,
		Storage: storage,
	})

	log.Println("JWKS successfully loaded from Keycloak")
	return nil
}

// Middleware to validate JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Auth Middleware: Checking request", r.Method, r.URL.Path)

		// Skip auth for the playground
		if r.URL.Path == "/playground" {
			next.ServeHTTP(w, r)
			return
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Unauthorized: Invalid Authorization format", authHeader)
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		log.Println("Auth Middleware: Token found, verifying...")

		// Parse and verify the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if keycloakJWKS == nil {
				return nil, fmt.Errorf("JWKS not initialized")
			}
			return keycloakJWKS.Keyfunc(token)
		}, jwt.WithLeeway(2*time.Minute))

		// Validations
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Print `iat` (issued at) claim
			if iat, ok := claims["iat"].(float64); ok {
				issuedAt := time.Unix(int64(iat), 0)
				serverTime := time.Now()

				log.Println("Token Issued At (iat):", issuedAt.UTC())
				log.Println("Server Current Time:", serverTime.UTC())
				log.Println("Time Difference:", serverTime.Sub(issuedAt))
			} else {
				log.Println("Error: Token missing `iat` claim")
			}
		} else {
			log.Println("Error: Unable to parse token claims")
		}
		if err != nil {
			log.Println("Unauthorized: Token validation error:", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
	
		if !token.Valid {
			log.Println("Unauthorized: Token is invalid")
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("Auth Middleware enabled: Token verified successfully")

		// Token is valid, pass request to next handler
		next.ServeHTTP(w, r)
	})
}
