package interpreter

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// BashCommand represents a bash command to be executed
type BashCommand struct {
	Command     string            `yaml:"command"`
	WorkingDir  string            `yaml:"working_dir,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Timeout     int               `yaml:"timeout,omitempty"` // timeout in seconds
}

// BashInterpreter implements the Interpreter interface for bash commands
type BashInterpreter struct{}

// NewBashInterpreter creates a new bash interpreter
func NewBashInterpreter() *BashInterpreter {
	return &BashInterpreter{}
}

// CanHandle returns true if the command type is "ktr@bash"
func (i *BashInterpreter) CanHandle(cmdType string) bool {
	return cmdType == "ktr@bash"
}

// Execute runs the bash command and returns the result
func (i *BashInterpreter) Execute(ctx context.Context, cmd Command, taskCtx *TaskContext) (*Result, error) {
	content, ok := cmd.Content.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid command content: must be a map")
	}

	// Extract command field
	cmdStr, ok := content["command"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid command content: command field is required and must be a string")
	}

	// Perform variable substitution in command
	substitutedCmd, err := taskCtx.Substitute(cmdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute variables in command: %w", err)
	}

	// Create the shell command
	shellCmd := exec.CommandContext(ctx, "bash", "-c", substitutedCmd)

	// Set environment variables
	shellCmd.Env = make([]string, 0)

	// Add all variables from context to environment
	for k, v := range taskCtx.Vars.Environment {
		shellCmd.Env = append(shellCmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range taskCtx.Vars.Secrets {
		shellCmd.Env = append(shellCmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Execute the command
	output, err := shellCmd.CombinedOutput()
	if err != nil {
		return &Result{
			Success: false,
			Output:  string(output),
			Error:   err,
		}, nil
	}

	return &Result{
		Success: true,
		Output:  strings.TrimSpace(string(output)),
	}, nil
}
