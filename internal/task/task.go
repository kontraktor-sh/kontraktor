package task

import (
	"context"
	"fmt"

	"github.com/kontraktor-sh/kontraktor/internal/output"
	"github.com/kontraktor-sh/kontraktor/internal/secret"
	"github.com/kontraktor-sh/kontraktor/internal/taskfile/interpreter"
	"github.com/kontraktor-sh/kontraktor/internal/vars"
)

// Task represents a single task in the taskfile
type Task struct {
	Desc        string                `yaml:"desc"`
	Args        []TaskArg             `yaml:"args,omitempty"`
	Cmds        []interpreter.Command `yaml:"cmds"`
	Environment map[string]string     `yaml:"environment,omitempty"`
}

// TaskArg represents a task argument
type TaskArg struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required,omitempty"`
	Default  string `yaml:"default,omitempty"`
}

// Executor handles task execution
type Executor struct {
	outputHandler *output.Handler
	secretManager *secret.Manager
	interpreter   interpreter.Interpreter
}

// NewExecutor creates a new task executor
func NewExecutor(outputHandler *output.Handler, secretManager *secret.Manager, interpreter interpreter.Interpreter) *Executor {
	return &Executor{
		outputHandler: outputHandler,
		secretManager: secretManager,
		interpreter:   interpreter,
	}
}

// Execute runs a task with the given arguments
func (e *Executor) Execute(ctx context.Context, task *Task, args map[string]interface{}) error {
	e.outputHandler.Debug("Executing task: %s", task.Desc)

	// Validate required arguments
	for _, arg := range task.Args {
		if arg.Required {
			if _, ok := args[arg.Name]; !ok {
				return fmt.Errorf("required argument '%s' not provided", arg.Name)
			}
		}
	}

	// Create task context
	taskCtx := &interpreter.TaskContext{
		Vars: vars.NewContext(),
	}
	taskCtx.Vars.Environment = task.Environment
	taskCtx.Vars.Args = args
	taskCtx.Vars.Secrets = make(map[string]string)

	// Load secrets
	if e.secretManager != nil {
		e.outputHandler.Debug("Loading secrets from vaults")
		secrets, err := e.secretManager.GetSecrets(ctx)
		if err != nil {
			return fmt.Errorf("failed to load secrets: %w", err)
		}
		taskCtx.Vars.Secrets = secrets
	}

	// Execute commands
	for _, cmd := range task.Cmds {
		content, _ := cmd.Content.(map[string]interface{})
		e.outputHandler.PrintCommand(cmd.Type, content)

		result, err := e.interpreter.Execute(ctx, cmd, taskCtx)

		if err != nil {
			e.outputHandler.Error("Command execution failed: %v", err)
			return fmt.Errorf("command execution failed: %w", err)
		}

		e.outputHandler.PrintResult(result)
		if !result.Success {
			return fmt.Errorf("command failed")
		}
	}

	e.outputHandler.Info("Task completed successfully")
	return nil
}
