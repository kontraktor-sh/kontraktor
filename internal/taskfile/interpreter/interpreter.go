package interpreter

import (
	"context"
	"fmt"

	"github.com/kontraktor-sh/kontraktor/internal/vars"
)

// TaskContext holds the execution context for a task
type TaskContext struct {
	Vars     *vars.Context
	TaskName string
}

// Command represents a command to be executed by an interpreter
type Command struct {
	Type    string
	Content interface{}
}

// Result represents the result of a command execution
type Result struct {
	Success bool
	Output  string
	Error   error
}

// Interpreter defines the interface that all command interpreters must implement
type Interpreter interface {
	// Execute runs the command and returns a result
	Execute(ctx context.Context, cmd Command, taskCtx *TaskContext) (*Result, error)

	// CanHandle returns true if this interpreter can handle the given command type
	CanHandle(cmdType string) bool
}

// Registry holds all available interpreters
type Registry struct {
	interpreters []Interpreter
}

// NewRegistry creates a new interpreter registry
func NewRegistry() *Registry {
	return &Registry{
		interpreters: make([]Interpreter, 0),
	}
}

// Register adds a new interpreter to the registry
func (r *Registry) Register(interpreter Interpreter) {
	r.interpreters = append(r.interpreters, interpreter)
}

// GetInterpreter returns the appropriate interpreter for the given command type
func (r *Registry) GetInterpreter(cmdType string) (Interpreter, error) {
	for _, interpreter := range r.interpreters {
		if interpreter.CanHandle(cmdType) {
			return interpreter, nil
		}
	}
	return nil, fmt.Errorf("no interpreter found for command type: %s", cmdType)
}

// NewTaskContext creates a new task context
func NewTaskContext() *TaskContext {
	return &TaskContext{
		Vars: vars.NewContext(),
	}
}

// SetEnvironment sets environment variables
func (c *TaskContext) SetEnvironment(env map[string]string) {
	for k, v := range env {
		c.Vars.Environment[k] = v
	}
}

// SetSecrets sets vault secrets
func (c *TaskContext) SetSecrets(secrets map[string]string) {
	for k, v := range secrets {
		c.Vars.Secrets[k] = v
	}
}

// SetArgs sets task arguments
func (c *TaskContext) SetArgs(args map[string]interface{}) {
	for k, v := range args {
		c.Vars.Args[k] = v
	}
}

// Substitute performs variable substitution in a string
func (c *TaskContext) Substitute(input string) (string, error) {
	return c.Vars.Substitutor.Substitute(input, c.Vars)
}
