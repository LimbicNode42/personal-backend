package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
	"net/http"
	"strings"
	"crypto/tls"

	"github.com/go-resty/resty/v2"
)

const (
	KeycloakURL  = "https://192.168.0.109:8443"
	Realm        = "shadow"
	ClientID     = "api-dev-site"
)

// Will need to be changed to retrieve from Infiscal
var ClientSecret = os.Getenv("KEYCLOAK_DEV_SITE_API")

// Middleware to validate JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate Secret
		if !ValidateSecret(tokenString) {
			log.Println("❌ Unauthorized: Invalid token")
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// // Validate Token
		// token, err := ValidateJWT(tokenString)
		// if err != nil || !token.Valid {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// 	c.Abort()
		// 	return
		// }

		// Token is valid, pass request to next handler
		next.ServeHTTP(w, r)
	})
}

// Validate Secret
func ValidateSecret(token string) bool {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // TODO: MUST be removed before deployed to production
	resp, err := client.R().
		SetBasicAuth(ClientID, ClientSecret).
		SetFormData(map[string]string{
			"token": token,
		}).
		Post(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", KeycloakURL, Realm))

	if err != nil {
		log.Println("❌ Error calling Keycloak introspection endpoint:", err)
		return false
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Println("❌ Error parsing Keycloak response:", err)
		return false
	}

	// Check if the token is active
	return result["active"].(bool)
}

// // ValidateJWT checks the JWT token against Keycloak public keys
// func ValidateJWT(tokenString string) (*jwt.Token, error) {
// 	// Get Keycloak public key dynamically
// 	cert, err := GetKeycloakPublicKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return jwt.ParseRSAPublicKeyFromPEM(cert)
// 	})
// }

// // Fetch Keycloak Public Key
// func GetKeycloakPublicKey() ([]byte, error) {
// 	client := resty.New()
// 	resp, err := client.R().Get(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", KeycloakURL, Realm))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var keyData map[string]interface{}
// 	if err := json.Unmarshal(resp.Body(), &keyData); err != nil {
// 		return nil, err
// 	}

// 	keys := keyData["keys"].([]interface{})
// 	kid := keys[0].(map[string]interface{})["x5c"].([]interface{})[0].(string)

// 	// Convert to PEM format
// 	cert := "-----BEGIN CERTIFICATE-----\n" + kid + "\n-----END CERTIFICATE-----"
// 	return []byte(cert), nil
// }
