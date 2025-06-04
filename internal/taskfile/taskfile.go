// Package taskfile provides structures for Kontraktor taskfiles.
package taskfile

// Task represents a single task in the taskfile
// desc: description
// args: list of task arguments
// cmds: list of shell commands
type Task struct {
	Desc        string            `yaml:"desc"`
	Args        []TaskArg         `yaml:"args,omitempty"`
	Cmds        []TaskCmd         `yaml:"cmds"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

// Vaults represents supported secret vaults configuration
// Currently only Azure Key Vault is supported
type Vaults struct {
	AzureKeyVault map[string]AzureKeyVaultConfig `yaml:"azure_keyvault,omitempty"`
}

// AzureKeyVaultConfig holds the config for a single Azure Key Vault
// keyvault_name: the name of the Azure Key Vault
// secrets: map of environment variable name to secret name in Key Vault
type AzureKeyVaultConfig struct {
	KeyVaultName string            `yaml:"keyvault_name"`
	Secrets      map[string]string `yaml:"secrets"`
}

// Taskfile represents the root of a taskfile.ktr.yml
// version: 0.3
// tasks: map of task name to Task
type Taskfile struct {
	Version     string            `yaml:"version"`
	Imports     []string          `yaml:"imports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Vaults      *Vaults           `yaml:"vaults,omitempty"`
	Tasks       map[string]Task   `yaml:"tasks"`
}
