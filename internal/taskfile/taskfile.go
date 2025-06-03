// Package taskfile provides structures and parsing for Kontraktor taskfiles.
package taskfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
)

// Taskfile represents the root of a taskfile.ktr.yml
// version: 0.3
// tasks: map of task name to Task
//
type Taskfile struct {
	Version     string            `yaml:"version"`
	Imports     []string          `yaml:"imports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Tasks       map[string]Task   `yaml:"tasks"`
}

// TaskCmd represents either a shell command or a task reference
// If Task is non-empty, it's a task reference; otherwise, use Cmd
//
type TaskCmd struct {
	Cmd  string
	Task string
	Uses string
	Args map[string]interface{}
}

// UnmarshalYAML implements custom YAML unmarshalling for TaskCmd
func (t *TaskCmd) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		t.Cmd = value.Value
		return nil
	}
	if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content); i += 2 {
			k := value.Content[i]
			v := value.Content[i+1]
			if k.Value == "task" {
				t.Task = v.Value
			}
			if k.Value == "uses" {
				t.Uses = v.Value
			}
			if k.Value == "args" {
				var m map[string]interface{}
				if err := v.Decode(&m); err != nil {
					return err
				}
				t.Args = m
			}
		}
		return nil
	}
	return fmt.Errorf("invalid cmd entry: %v", value)
}

// TaskArg represents an argument for a task
// name: argument name
// type: argument type (string, [], bool, number)
// default: default value (interface{})
type TaskArg struct {
	Name    string      `yaml:"name"`
	Type    string      `yaml:"type"`
	Default interface{} `yaml:"default,omitempty"`
}

// Task represents a single task in the taskfile
// desc: description
// args: list of task arguments
// cmds: list of shell commands
//
type Task struct {
	Desc        string            `yaml:"desc"`
	Args        []TaskArg         `yaml:"args,omitempty"`
	Cmds        []TaskCmd         `yaml:"cmds"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

// ParseTaskfile reads and parses a taskfile.ktr.yml from the given path, recursively loading imports.
func ParseTaskfile(path string) (*Taskfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open taskfile: %w", err)
	}
	defer f.Close()

	var tf Taskfile
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&tf); err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	// Recursively load imports
	for _, importPath := range tf.Imports {
		var importFile string
		if isHTTPImport(importPath) {
			importFile, err = downloadToTemp(importPath)
			if err != nil {
				return nil, fmt.Errorf("download import %s: %w", importPath, err)
			}
		} else if isGitImport(importPath) {
			importFile, err = cloneAndGetFile(importPath)
			if err != nil {
				return nil, fmt.Errorf("git import %s: %w", importPath, err)
			}
		} else {
			importFile = importPath
		}
		imported, err := ParseTaskfile(importFile)
		if err != nil {
			return nil, fmt.Errorf("import %s: %w", importPath, err)
		}
		// Merge imported tasks, but do not override main file tasks
		for k, v := range imported.Tasks {
			if _, exists := tf.Tasks[k]; !exists {
				tf.Tasks[k] = v
			}
		}
	}

	return &tf, nil
}

func isHTTPImport(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func isGitImport(path string) bool {
	return strings.Contains(path, ".git//")
}

func downloadToTemp(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "ktr-import-*.yml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func cloneAndGetFile(gitImport string) (string, error) {
	parts := strings.SplitN(gitImport, ".git//", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid git import: %s", gitImport)
	}
	repoURL := parts[0] + ".git"
	fileInRepo := parts[1]

	tmpDir, err := os.MkdirTemp("", "ktr-git-*")
	if err != nil {
		return "", err
	}
	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		return "", err
	}
	return filepath.Join(tmpDir, fileInRepo), nil
} 