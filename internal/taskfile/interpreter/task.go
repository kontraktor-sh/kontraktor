package interpreter

import (
	"context"
	"fmt"

	"github.com/kontraktor-sh/kontraktor/internal/vars"
)

// TaskCommand represents a reference to another task
type TaskCommand struct {
	Name string
	Args map[string]interface{}
}

// TaskInterpreter implements the Interpreter interface for task references
type TaskInterpreter struct {
	taskExecutor func(ctx context.Context, taskName string, args map[string]interface{}, taskCtx *TaskContext) (*Result, error)
}

// NewTaskInterpreter creates a new task interpreter
func NewTaskInterpreter(executor func(ctx context.Context, taskName string, args map[string]interface{}, taskCtx *TaskContext) (*Result, error)) *TaskInterpreter {
	return &TaskInterpreter{
		taskExecutor: executor,
	}
}

// CanHandle returns true if the command type is "task"
func (i *TaskInterpreter) CanHandle(cmdType string) bool {
	return cmdType == "task"
}

// Execute runs the referenced task and returns the result
func (i *TaskInterpreter) Execute(ctx context.Context, cmd Command, taskCtx *TaskContext) (*Result, error) {
	taskCmd, ok := cmd.Content.(TaskCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command content for task interpreter")
	}

	if taskCmd.Name == "" {
		return nil, fmt.Errorf("task name is required")
	}

	// Merge task arguments with the current context
	mergedArgs := make(map[string]interface{})
	for k, v := range taskCtx.Vars.Args {
		mergedArgs[k] = v
	}
	for k, v := range taskCmd.Args {
		mergedArgs[k] = v
	}

	// Create a new task context for the referenced task
	newTaskCtx := &TaskContext{
		Vars:     vars.NewContext(),
		TaskName: taskCmd.Name,
	}

	// Copy environment variables and secrets
	for k, v := range taskCtx.Vars.Environment {
		newTaskCtx.Vars.Environment[k] = v
	}
	for k, v := range taskCtx.Vars.Secrets {
		newTaskCtx.Vars.Secrets[k] = v
	}
	newTaskCtx.Vars.Args = mergedArgs

	// Execute the referenced task
	return i.taskExecutor(ctx, taskCmd.Name, mergedArgs, newTaskCtx)
}
