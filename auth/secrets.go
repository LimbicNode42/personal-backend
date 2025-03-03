package auth

import (
	"os"
	"fmt"
	"context"

	infisical "github.com/infisical/go-sdk"
)


var InfURL = "http://192.168.0.108:8080"
var InfProjectID = os.Getenv("INF_DEV_PROJECT_ID")
var InfClientID = os.Getenv("INF_DEV_API_CLIENT_ID")
var InfClientSecret = os.Getenv("INF_DEV_API_CLIENT_SECRET")

func infisicalLogin() infisical.InfisicalClientInterface {
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl: InfURL,
    	AutoTokenRefresh: true,
	})

	_, err := client.Auth().UniversalAuthLogin(InfClientID, InfClientSecret)

	if err != nil {
		fmt.Printf("Authentication failed: %v", err)
		os.Exit(1)
	}

	return client
}

func InfisicalGetSecrets(path string) []infisical.Secret {
	client := infisicalLogin()
	secrets, err := client.Secrets().List(infisical.ListSecretsOptions{
		// SecretKey:   "API_KEY",
		ProjectID:   InfProjectID,
		ProjectSlug: "dev-site",
		Environment: "prod",
		SecretPath:  path,
	})

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	return secrets
}
