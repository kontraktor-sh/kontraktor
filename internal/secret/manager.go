package secret

import (
	"context"
	"fmt"
)

// Manager handles secret management
type Manager struct {
	vaults map[string]Vault
}

// Vault represents a secret vault
type Vault interface {
	GetSecrets(ctx context.Context) (map[string]string, error)
}

// NewManager creates a new secret manager
func NewManager() *Manager {
	return &Manager{
		vaults: make(map[string]Vault),
	}
}

// RegisterVault registers a new vault
func (m *Manager) RegisterVault(name string, vault Vault) {
	m.vaults[name] = vault
}

// GetSecrets retrieves secrets from all registered vaults
func (m *Manager) GetSecrets(ctx context.Context) (map[string]string, error) {
	secrets := make(map[string]string)

	for name, vault := range m.vaults {
		vaultSecrets, err := vault.GetSecrets(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get secrets from vault %s: %w", name, err)
		}

		// Merge secrets, with later vaults taking precedence
		for k, v := range vaultSecrets {
			secrets[k] = v
		}
	}

	return secrets, nil
}
