package interpreter

import "context"

// NewDefaultRegistry creates a new registry with all default interpreters
func NewDefaultRegistry(taskExecutor func(ctx context.Context, taskName string, args map[string]interface{}, taskCtx *TaskContext) (*Result, error)) *Registry {
	registry := NewRegistry()

	// Register default interpreters
	registry.Register(NewBashInterpreter())
	registry.Register(NewPythonInterpreter())
	registry.Register(NewDockerInterpreter())
	registry.Register(NewTaskInterpreter(taskExecutor))

	return registry
}
