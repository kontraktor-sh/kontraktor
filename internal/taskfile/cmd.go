// cmd.go
// Defines TaskCmd and its YAML unmarshalling for Kontraktor taskfiles.
package taskfile

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// TaskCmd represents a command that can be executed by an interpreter
type TaskCmd struct {
	Type    string                 `yaml:"type"`
	Content map[string]interface{} `yaml:"content"`
}

// UnmarshalYAML implements custom YAML unmarshalling for TaskCmd
func (t *TaskCmd) UnmarshalYAML(value *yaml.Node) error {
	// Handle simple string commands (backward compatibility)
	if value.Kind == yaml.ScalarNode {
		t.Type = "bash"
		t.Content = map[string]interface{}{
			"command": value.Value,
		}
		return nil
	}

	// Handle structured commands
	if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content); i += 2 {
			k := value.Content[i]
			v := value.Content[i+1]

			switch k.Value {
			case "type":
				t.Type = v.Value
			case "content":
				var m map[string]interface{}
				if err := v.Decode(&m); err != nil {
					return err
				}
				t.Content = m
			}
		}
		return nil
	}

	return fmt.Errorf("invalid cmd entry: %v", value)
}
