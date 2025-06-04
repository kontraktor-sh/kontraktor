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
		// First, unmarshal into a map to check the command type
		var cmdMap map[string]interface{}
		if err := value.Decode(&cmdMap); err != nil {
			return err
		}

		// Check for task reference
		if taskName, ok := cmdMap["task"].(string); ok {
			t.Type = "task"
			t.Content = map[string]interface{}{
				"name": taskName,
			}
			return nil
		}

		// Check for uses command
		if _, ok := cmdMap["uses"].(string); ok {
			t.Type = "uses"
			t.Content = cmdMap
			return nil
		}

		// Handle explicit type and content
		if cmdType, ok := cmdMap["type"].(string); ok {
			t.Type = cmdType
			if content, ok := cmdMap["content"].(map[string]interface{}); ok {
				t.Content = content
			} else {
				t.Content = cmdMap
			}
			return nil
		}

		// If no specific type is found, treat as bash command
		t.Type = "bash"
		t.Content = cmdMap
		return nil
	}

	return fmt.Errorf("invalid cmd entry: %v", value)
}
