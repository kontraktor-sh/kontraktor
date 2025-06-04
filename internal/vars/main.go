package vars

import (
	"fmt"
	"regexp"
)

// VariableType represents the type of variable
type VariableType string

const (
	TypeEnv     VariableType = "env"     // Environment variables
	TypeSecret  VariableType = "secret"  // Vault secrets
	TypeArg     VariableType = "arg"     // Task arguments
	TypeUnknown VariableType = "unknown" // Unknown type
)

// Variable represents a variable with its type and value
type Variable struct {
	Type  VariableType
	Name  string
	Value string
}

// Context holds all available variables for substitution
type Context struct {
	Environment map[string]string      // Environment variables
	Secrets     map[string]string      // Vault secrets
	Args        map[string]interface{} // Task arguments
	Substitutor *Substitutor           // Variable substitutor
}

// NewContext creates a new variable context
func NewContext() *Context {
	return &Context{
		Environment: make(map[string]string),
		Secrets:     make(map[string]string),
		Args:        make(map[string]interface{}),
		Substitutor: NewSubstitutor(),
	}
}

// GetVariable retrieves a variable by name, checking all sources
func (c *Context) GetVariable(name string) (*Variable, error) {
	// Check environment variables
	if value, ok := c.Environment[name]; ok {
		return &Variable{Type: TypeEnv, Name: name, Value: value}, nil
	}

	// Check secrets
	if value, ok := c.Secrets[name]; ok {
		return &Variable{Type: TypeSecret, Name: name, Value: value}, nil
	}

	// Check arguments
	if value, ok := c.Args[name]; ok {
		// Convert argument value to string
		strValue := fmt.Sprintf("%v", value)
		return &Variable{Type: TypeArg, Name: name, Value: strValue}, nil
	}

	return nil, fmt.Errorf("variable '%s' not found", name)
}

// Substitutor handles variable substitution
type Substitutor struct {
	varRegex *regexp.Regexp
}

// NewSubstitutor creates a new variable substitutor
func NewSubstitutor() *Substitutor {
	return &Substitutor{
		varRegex: regexp.MustCompile(`\${([^}]+)}`),
	}
}

// Substitute performs variable substitution using the context
func (s *Substitutor) Substitute(input string, ctx *Context) (string, error) {
	return s.varRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := match[2 : len(match)-1]

		// Get variable from context
		variable, err := ctx.GetVariable(varName)
		if err != nil {
			return match // Return original if not found
		}

		// Perform nested substitution on the value
		substituted, err := s.Substitute(variable.Value, ctx)
		if err != nil {
			return match // Return original if substitution fails
		}

		return substituted
	}), nil
}

// SubstituteMap performs variable substitution on a map of strings
func (s *Substitutor) SubstituteMap(input map[string]string, ctx *Context) (map[string]string, error) {
	result := make(map[string]string)
	for k, v := range input {
		substituted, err := s.Substitute(v, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to substitute in value for '%s': %w", k, err)
		}
		result[k] = substituted
	}
	return result, nil
}
