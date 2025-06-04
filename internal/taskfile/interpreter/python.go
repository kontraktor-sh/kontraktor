package interpreter

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// PythonCommand represents a Python command to be executed
type PythonCommand struct {
	Script string
	Args   []string
}

// PythonInterpreter implements the Interpreter interface for Python commands
type PythonInterpreter struct{}

// NewPythonInterpreter creates a new Python interpreter
func NewPythonInterpreter() *PythonInterpreter {
	return &PythonInterpreter{}
}

// CanHandle returns true if the command type is "python"
func (i *PythonInterpreter) CanHandle(cmdType string) bool {
	return cmdType == "python"
}

// Execute runs the Python command and returns the result
func (i *PythonInterpreter) Execute(ctx context.Context, cmd Command, taskCtx *TaskContext) (*Result, error) {
	pythonCmd, ok := cmd.Content.(PythonCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command content for python interpreter")
	}

	// Create the Python command
	args := append([]string{"-c", pythonCmd.Script}, pythonCmd.Args...)
	shellCmd := exec.CommandContext(ctx, "python3", args...)

	// Set environment variables
	shellCmd.Env = make([]string, 0)
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
