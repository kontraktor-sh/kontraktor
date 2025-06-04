// Package taskfile provides structures for Kontraktor taskfiles.
package taskfile

import (
	"fmt"

	"github.com/kontraktor-sh/kontraktor/internal/env"
	"github.com/kontraktor-sh/kontraktor/internal/vars"
)

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

// Validate validates the taskfile
func (tf *Taskfile) Validate() error {
	validator := env.NewValidator()

	// Validate global environment variables
	if err := validator.ValidateMap(tf.Environment); err != nil {
		return fmt.Errorf("invalid global environment: %w", err)
	}

	// Validate task environment variables
	for taskName, task := range tf.Tasks {
		if err := validator.ValidateMap(task.Environment); err != nil {
			return fmt.Errorf("invalid environment in task '%s': %w", taskName, err)
		}
	}

	return nil
}

// ProcessVariables performs variable substitution in all variable sources
func (tf *Taskfile) ProcessVariables() error {
	ctx := vars.NewContext()

	// Add global environment variables to context
	for k, v := range tf.Environment {
		ctx.Environment[k] = v
	}

	// Process global environment variables
	processed, err := ctx.Substitutor.SubstituteMap(tf.Environment, ctx)
	if err != nil {
		return fmt.Errorf("failed to process global environment: %w", err)
	}
	tf.Environment = processed

	// Process task environment variables
	for taskName, task := range tf.Tasks {
		if task.Environment != nil {
			// Create task-specific context
			taskCtx := vars.NewContext()

			// Add global environment to task context
			for k, v := range tf.Environment {
				taskCtx.Environment[k] = v
			}

			// Process task environment variables
			processed, err := ctx.Substitutor.SubstituteMap(task.Environment, taskCtx)
			if err != nil {
				return fmt.Errorf("failed to process environment in task '%s': %w", taskName, err)
			}
			task.Environment = processed
			tf.Tasks[taskName] = task
		}
	}

	return nil
}
