package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rafaelherik/kontraktor-sh/internal/taskfile"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestExecuteTaskWithArgs(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "kontraktor-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test taskfile
	testTaskfile := &taskfile.Taskfile{
		Version: "0.3",
		Environment: map[string]string{
			"GLOBAL_VAR": "global-value",
		},
		Tasks: map[string]taskfile.Task{
			"test-task": {
				Desc: "Test task with args",
				Args: []taskfile.TaskArg{
					{
						Name:    "name",
						Type:    "string",
						Default: "default-name",
					},
				},
				Environment: map[string]string{
					"TASK_VAR": "task-value",
				},
				Cmds: []taskfile.TaskCmd{
					{Cmd: "echo 'Hello, ${name}!'"},
				},
			},
		},
	}

	// Write test taskfile to disk
	taskfilePath := filepath.Join(tmpDir, "test.ktr.yml")
	data, err := yaml.Marshal(testTaskfile)
	assert.NoError(t, err)
	err = os.WriteFile(taskfilePath, data, 0644)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		taskName      string
		args          map[string]string
		expectedError bool
	}{
		{
			name:          "execute with default args",
			taskName:      "test-task",
			args:          map[string]string{},
			expectedError: false,
		},
		{
			name:     "execute with custom args",
			taskName: "test-task",
			args: map[string]string{
				"name": "custom-name",
			},
			expectedError: false,
		},
		{
			name:          "non-existent task",
			taskName:      "non-existent",
			args:          map[string]string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := make(map[string]bool)
			err := executeTaskWithArgs(tt.taskName, testTaskfile, visited, tt.args, nil, []int{}, nil)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunShellCommand(t *testing.T) {
	tests := []struct {
		name          string
		cmd           string
		env           map[string]string
		expectedError bool
	}{
		{
			name:          "successful command",
			cmd:           "echo 'test'",
			env:           nil,
			expectedError: false,
		},
		{
			name:          "command with environment",
			cmd:           "echo $TEST_VAR",
			env:           map[string]string{"TEST_VAR": "test-value"},
			expectedError: false,
		},
		{
			name:          "failing command",
			cmd:           "exit 1",
			env:           nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runShellCommand(tt.cmd, tt.env)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
