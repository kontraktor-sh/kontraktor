package vault

import (
	"context"
	"testing"

	"github.com/kontraktor-sh/kontraktor/internal/taskfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAzureClient is a mock implementation of the AzureClient interface
type MockAzureClient struct {
	mock.Mock
}

func (m *MockAzureClient) GetSecret(ctx context.Context, vaultName, secretName string) (string, error) {
	args := m.Called(ctx, vaultName, secretName)
	return args.String(0), args.Error(1)
}

func TestFetchAzureSecrets(t *testing.T) {
	tests := []struct {
		name        string
		config      taskfile.AzureKeyVaultConfig
		mockSecrets map[string]string
		want        map[string]string
		wantErr     bool
	}{
		{
			name: "successful secret fetch",
			config: taskfile.AzureKeyVaultConfig{
				KeyVaultName: "test-vault",
				Secrets: map[string]string{
					"API_KEY":     "api-secret",
					"DB_PASSWORD": "db-secret",
				},
			},
			mockSecrets: map[string]string{
				"api-secret": "actual-api-key",
				"db-secret":  "actual-db-password",
			},
			want: map[string]string{
				"API_KEY":     "actual-api-key",
				"DB_PASSWORD": "actual-db-password",
			},
			wantErr: false,
		},
		{
			name: "missing secret",
			config: taskfile.AzureKeyVaultConfig{
				KeyVaultName: "test-vault",
				Secrets: map[string]string{
					"API_KEY": "missing-secret",
				},
			},
			mockSecrets: map[string]string{},
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAzureClient)

			// Set up mock expectations
			for secretName, secretValue := range tt.mockSecrets {
				mockClient.On("GetSecret", mock.Anything, tt.config.KeyVaultName, secretName).
					Return(secretValue, nil)
			}

			// For missing secrets, return an error
			for _, secretName := range tt.config.Secrets {
				if _, exists := tt.mockSecrets[secretName]; !exists {
					mockClient.On("GetSecret", mock.Anything, tt.config.KeyVaultName, secretName).
						Return("", assert.AnError)
				}
			}

			got, err := FetchSecrets(context.Background(), mockClient, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
