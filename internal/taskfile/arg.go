// arg.go
// Defines TaskArg for task arguments in Kontraktor taskfiles.
package taskfile

// TaskArg represents an argument for a task.
// name: argument name
// type: argument type (string, [], bool, number)
// default: default value (interface{})
type TaskArg struct {
	Name    string      `yaml:"name"`
	Type    string      `yaml:"type"`
	Default interface{} `yaml:"default,omitempty"`
}
