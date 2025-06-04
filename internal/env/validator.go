package env

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents an environment variable validation error
type ValidationError struct {
	VarName string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid environment variable '%s': %s", e.VarName, e.Reason)
}

// Validator handles environment variable validation
type Validator struct {
	// Regular expression for valid environment variable names
	// Follows POSIX standard: [a-zA-Z_][a-zA-Z0-9_]*
	nameRegex *regexp.Regexp
	// List of reserved environment variable names
	reservedNames map[string]bool
}

// NewValidator creates a new environment variable validator
func NewValidator() *Validator {
	return &Validator{
		nameRegex: regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`),
		reservedNames: map[string]bool{
			"PATH":    true,
			"HOME":    true,
			"USER":    true,
			"SHELL":   true,
			"PWD":     true,
			"OLDPWD":  true,
			"TERM":    true,
			"LANG":    true,
			"LC_ALL":  true,
			"TZ":      true,
			"EDITOR":  true,
			"VISUAL":  true,
			"PAGER":   true,
			"MANPATH": true,
		},
	}
}

// ValidateName validates an environment variable name
func (v *Validator) ValidateName(name string) error {
	if name == "" {
		return &ValidationError{VarName: name, Reason: "name cannot be empty"}
	}

	if !v.nameRegex.MatchString(name) {
		return &ValidationError{
			VarName: name,
			Reason:  "name must start with a letter or underscore and contain only letters, numbers, and underscores",
		}
	}

	if v.reservedNames[strings.ToUpper(name)] {
		return &ValidationError{
			VarName: name,
			Reason:  "name is reserved and cannot be overridden",
		}
	}

	return nil
}

// ValidateValue validates an environment variable value
func (v *Validator) ValidateValue(value string) error {
	// Check for null bytes
	if strings.ContainsRune(value, 0) {
		return &ValidationError{
			VarName: "",
			Reason:  "value cannot contain null bytes",
		}
	}

	return nil
}

// ValidateMap validates a map of environment variables
func (v *Validator) ValidateMap(env map[string]string) error {
	for name, value := range env {
		if err := v.ValidateName(name); err != nil {
			return err
		}
		if err := v.ValidateValue(value); err != nil {
			return err
		}
	}
	return nil
}
