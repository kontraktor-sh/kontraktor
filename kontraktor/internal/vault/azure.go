// Package vault provides secret retrieval from Azure Key Vault.
package vault

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/rafaelherik/kontraktor-sh/internal/taskfile"
)

// AzureClient defines the interface for Azure Key Vault operations
type AzureClient interface {
	GetSecret(ctx context.Context, vaultName, secretName string) (string, error)
}

// AzureKeyVaultClient implements AzureClient using the Azure SDK
type AzureKeyVaultClient struct {
	client *azsecrets.Client
}

// NewAzureKeyVaultClient creates a new AzureKeyVaultClient
func NewAzureKeyVaultClient(vaultName string) (*AzureKeyVaultClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure credential: %w", err)
	}
	client, err := azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", vaultName), cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Key Vault client: %w", err)
	}
	return &AzureKeyVaultClient{client: client}, nil
}

// GetSecret retrieves a secret from Azure Key Vault
func (c *AzureKeyVaultClient) GetSecret(ctx context.Context, vaultName, secretName string) (string, error) {
	resp, err := c.client.GetSecret(ctx, secretName, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}
	return *resp.Value, nil
}

// FetchAzureSecrets fetches secrets from Azure Key Vault as specified in the config.
// Returns a map of environment variable names to secret values.
func FetchAzureSecrets(ctx context.Context, config taskfile.AzureKeyVaultConfig) (map[string]string, error) {
	client, err := NewAzureKeyVaultClient(config.KeyVaultName)
	if err != nil {
		return nil, err
	}
	return FetchSecrets(ctx, client, config)
}

// FetchSecrets fetches secrets using the provided AzureClient
func FetchSecrets(ctx context.Context, client AzureClient, config taskfile.AzureKeyVaultConfig) (map[string]string, error) {
	result := make(map[string]string)
	for envVar, secretName := range config.Secrets {
		value, err := client.GetSecret(ctx, config.KeyVaultName, secretName)
		if err != nil {
			return nil, err
		}
		result[envVar] = value
	}
	return result, nil
}
