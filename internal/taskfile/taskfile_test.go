package taskfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParseTaskfile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "kontraktor-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test taskfile
	testTaskfile := Taskfile{
		Version: "0.3",
		Environment: map[string]string{
			"FOO": "bar",
		},
		Tasks: map[string]Task{
			"hello": {
				Desc: "A simple hello world task",
				Cmds: []TaskCmd{
					{
						Cmd: "echo 'Hello, World!'",
					},
				},
			},
		},
	}

	// Write test taskfile to disk
	taskfilePath := filepath.Join(tmpDir, "test.ktr.yml")
	taskfileContent := `version: "0.3"
environment:
  FOO: bar
tasks:
  hello:
    desc: A simple hello world task
    cmds:
      - echo 'Hello, World!'
`
	err = os.WriteFile(taskfilePath, []byte(taskfileContent), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		path          string
		expectedError bool
	}{
		{
			name:          "valid taskfile",
			path:          taskfilePath,
			expectedError: false,
		},
		{
			name:          "non-existent file",
			path:          "non-existent.ktr.yml",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf, err := ParseTaskfile(tt.path)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, tf)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tf)
				assert.Equal(t, testTaskfile.Version, tf.Version)
				assert.Equal(t, testTaskfile.Environment, tf.Environment)
				// Compare tasks individually to handle nil vs empty map differences
				for name, expectedTask := range testTaskfile.Tasks {
					actualTask, exists := tf.Tasks[name]
					assert.True(t, exists, "Task %s should exist", name)
					assert.Equal(t, expectedTask.Desc, actualTask.Desc)
					assert.Equal(t, expectedTask.Args, actualTask.Args)
					assert.Equal(t, expectedTask.Environment, actualTask.Environment)
					assert.Equal(t, len(expectedTask.Cmds), len(actualTask.Cmds))
					for i, expectedCmd := range expectedTask.Cmds {
						assert.Equal(t, expectedCmd.Cmd, actualTask.Cmds[i].Cmd)
						assert.Equal(t, expectedCmd.Task, actualTask.Cmds[i].Task)
						assert.Equal(t, expectedCmd.Uses, actualTask.Cmds[i].Uses)
						assert.Equal(t, expectedCmd.Args, actualTask.Cmds[i].Args)
					}
				}
			}
		})
	}
}

func TestTaskCmdUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name          string
		yaml          string
		expectedCmd   TaskCmd
		expectedError bool
	}{
		{
			name: "simple command",
			yaml: "echo 'Hello, World!'",
			expectedCmd: TaskCmd{
				Cmd: "echo 'Hello, World!'",
			},
			expectedError: false,
		},
		{
			name: "task reference",
			yaml: "task: hello",
			expectedCmd: TaskCmd{
				Task: "hello",
			},
			expectedError: false,
		},
		{
			name: "uses with args",
			yaml: `
uses: docker-build
args:
  image: test-image
`,
			expectedCmd: TaskCmd{
				Uses: "docker-build",
				Args: map[string]interface{}{
					"image": "test-image",
				},
			},
			expectedError: false,
		},
		{
			name:          "invalid yaml",
			yaml:          "invalid: yaml: content",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TaskCmd
			err := yaml.Unmarshal([]byte(tt.yaml), &cmd)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmd, cmd)
			}
		})
	}
}
