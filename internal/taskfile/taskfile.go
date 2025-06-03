// Package taskfile provides structures and parsing for Kontraktor taskfiles.
package taskfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Taskfile represents the root of a taskfile.ktr.yml
// version: 0.3
// tasks: map of task name to Task
//
type Taskfile struct {
	Version string           `yaml:"version"`
	Tasks   map[string]Task  `yaml:"tasks"`
}

// Task represents a single task in the taskfile
// desc: description
// cmds: list of shell commands
//
type Task struct {
	Desc string   `yaml:"desc"`
	Cmds []string `yaml:"cmds"`
}

// ParseTaskfile reads and parses a taskfile.ktr.yml from the given path.
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
	return &tf, nil
} 