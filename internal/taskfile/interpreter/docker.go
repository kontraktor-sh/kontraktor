package interpreter

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// DockerCommand represents a Docker command to be executed
type DockerCommand struct {
	Image       string
	Command     []string
	Environment map[string]string
	Volumes     map[string]string
	Network     string
}

// DockerInterpreter implements the Interpreter interface for Docker commands
type DockerInterpreter struct{}

// NewDockerInterpreter creates a new Docker interpreter
func NewDockerInterpreter() *DockerInterpreter {
	return &DockerInterpreter{}
}

// CanHandle returns true if the command type is "docker"
func (i *DockerInterpreter) CanHandle(cmdType string) bool {
	return cmdType == "docker"
}

// Execute runs the Docker command and returns the result
func (i *DockerInterpreter) Execute(ctx context.Context, cmd Command, taskCtx *TaskContext) (*Result, error) {
	dockerCmd, ok := cmd.Content.(DockerCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command content for docker interpreter")
	}

	if dockerCmd.Image == "" {
		return nil, fmt.Errorf("docker image is required")
	}

	// Build docker run command
	args := []string{"run", "--rm"}

	// Add environment variables
	for k, v := range taskCtx.Vars.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range taskCtx.Vars.Secrets {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range dockerCmd.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	// Add volumes
	for host, container := range dockerCmd.Volumes {
		args = append(args, "-v", fmt.Sprintf("%s:%s", host, container))
	}

	// Add network if specified
	if dockerCmd.Network != "" {
		args = append(args, "--network", dockerCmd.Network)
	}

	// Add image and command
	args = append(args, dockerCmd.Image)
	args = append(args, dockerCmd.Command...)

	// Create the docker command
	shellCmd := exec.CommandContext(ctx, "docker", args...)

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
