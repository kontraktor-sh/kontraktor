package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kontraktor-sh/kontraktor/internal/cli"
	"github.com/kontraktor-sh/kontraktor/internal/secret"
	"github.com/kontraktor-sh/kontraktor/internal/task"
	"github.com/kontraktor-sh/kontraktor/internal/taskfile"
	"github.com/kontraktor-sh/kontraktor/internal/taskfile/interpreter"
)

func main() {
	// Parse command line flags
	config, err := cli.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create output handler
	outputHandler, err := config.CreateOutputHandler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output handler: %v\n", err)
		os.Exit(1)
	}

	// Create secret manager
	secretManager := secret.NewManager()
	// TODO: Register vaults based on configuration

	// Create interpreter registry
	registry := interpreter.NewRegistry()
	registry.Register(interpreter.NewBashInterpreter())
	// TODO: Register other interpreters

	// Create interpreter
	interpreter := interpreter.NewBashInterpreter()

	// Create task executor
	executor := task.NewExecutor(outputHandler, secretManager, interpreter)

	// Load the taskfile
	taskfile, err := taskfile.ParseTaskfile("taskfile.ktr.yml")
	if err != nil {
		outputHandler.Error("Failed to load taskfile: %v", err)
		os.Exit(1)
	}

	// Get the task definition
	taskDef, exists := taskfile.Tasks[config.TaskName]
	if !exists {
		outputHandler.Error("Task '%s' not found in taskfile", config.TaskName)
		os.Exit(1)
	}

	// Convert taskfile.Task to task.Task
	taskToExecute := &task.Task{
		Desc:        taskDef.Desc,
		Args:        convertTaskArgs(taskDef.Args),
		Cmds:        convertTaskCmds(taskDef.Cmds),
		Environment: make(map[string]string),
	}

	// Copy global environment variables
	for k, v := range taskfile.Environment {
		taskToExecute.Environment[k] = v
	}

	// Copy task-specific environment variables (overriding globals)
	for k, v := range taskDef.Environment {
		taskToExecute.Environment[k] = v
	}

	// Execute the task
	ctx := context.Background()
	args := make(map[string]interface{})
	for k, v := range config.TaskArgs {
		args[k] = v
	}
	err = executor.Execute(ctx, taskToExecute, args)
	if err != nil {
		outputHandler.Error("Task execution failed: %v", err)
		os.Exit(1)
	}
}

func convertTaskArgs(args []taskfile.TaskArg) []task.TaskArg {
	result := make([]task.TaskArg, len(args))
	for i, arg := range args {
		result[i] = task.TaskArg{
			Name:    arg.Name,
			Default: fmt.Sprint(arg.Default),
		}
	}
	return result
}

func convertTaskCmds(cmds []taskfile.TaskCmd) []interpreter.Command {
	result := make([]interpreter.Command, len(cmds))
	for i, cmd := range cmds {
		result[i] = interpreter.Command{
			Type:    cmd.Type,
			Content: cmd.Content,
		}
	}
	return result
}
