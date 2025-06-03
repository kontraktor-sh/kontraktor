// Package vault provides secret retrieval from Azure Key Vault.
package vault

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/rafaelherik/kontraktor-sh/kontraktor/internal/taskfile"
)

// FetchAzureSecrets fetches secrets from Azure Key Vault as specified in the config.
// Returns a map of environment variable names to secret values.
func FetchAzureSecrets(ctx context.Context, config taskfile.AzureKeyVaultConfig) (map[string]string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure credential: %w", err)
	}
	client, err := azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", config.KeyVaultName), cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Key Vault client: %w", err)
	}
	result := make(map[string]string)
	for envVar, secretName := range config.Secrets {
		resp, err := client.GetSecret(ctx, secretName, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get secret %s: %w", secretName, err)
		}
		result[envVar] = *resp.Value
	}
	return result, nil
} 