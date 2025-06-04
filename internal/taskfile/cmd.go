// cmd.go
// Defines TaskCmd and its YAML unmarshalling for Kontraktor taskfiles.
package taskfile

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// TaskCmd represents either a shell command or a task reference.
// If Task is non-empty, it's a task reference; otherwise, use Cmd.
type TaskCmd struct {
	Cmd  string
	Task string
	Uses string
	Args map[string]interface{}
}

// UnmarshalYAML implements custom YAML unmarshalling for TaskCmd.
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
