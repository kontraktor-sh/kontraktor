package output

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/kontraktor-sh/kontraktor/internal/taskfile/interpreter"
)

// VerbosityLevel represents the level of output verbosity
type VerbosityLevel int

const (
	// LevelSilent suppresses all output except critical errors
	LevelSilent VerbosityLevel = iota
	// LevelError shows only errors
	LevelError
	// LevelInfo shows normal output and errors
	LevelInfo
	// LevelDebug shows all output including debug information
	LevelDebug
)

// Handler manages command output and sensitive data masking
type Handler struct {
	maskPatterns []*regexp.Regexp
	verbosity    VerbosityLevel
	out          io.Writer
	err          io.Writer
}

// NewHandler creates a new output handler
func NewHandler() *Handler {
	return &Handler{
		maskPatterns: make([]*regexp.Regexp, 0),
		verbosity:    LevelInfo, // Default to Info level
		out:          os.Stdout,
		err:          os.Stderr,
	}
}

// SetLevel sets the verbosity level
func (h *Handler) SetLevel(level VerbosityLevel) {
	h.verbosity = level
}

// SetOutput sets the output writer
func (h *Handler) SetOutput(out io.Writer) {
	h.out = out
}

// SetError sets the error writer
func (h *Handler) SetError(err io.Writer) {
	h.err = err
}

// AddMaskPattern adds a regex pattern for masking sensitive data
func (h *Handler) AddMaskPattern(pattern string) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Log error but don't fail
		fmt.Fprintf(h.err, "Warning: invalid mask pattern '%s': %v\n", pattern, err)
		return
	}
	h.maskPatterns = append(h.maskPatterns, re)
}

// MaskSensitiveData masks sensitive information in the output
func (h *Handler) MaskSensitiveData(output string) string {
	masked := output
	for _, pattern := range h.maskPatterns {
		masked = pattern.ReplaceAllString(masked, "[MASKED]")
	}
	return masked
}

// Debug prints debug information if verbosity level is DebugLevel
func (h *Handler) Debug(format string, args ...interface{}) {
	if h.verbosity >= LevelDebug {
		fmt.Fprintf(h.out, "[DEBUG] "+format+"\n", args...)
	}
}

// Info prints information if verbosity level is InfoLevel or higher
func (h *Handler) Info(format string, args ...interface{}) {
	if h.verbosity >= LevelInfo {
		fmt.Fprintf(h.out, format+"\n", args...)
	}
}

// Error prints error information if verbosity level is ErrorLevel or higher
func (h *Handler) Error(format string, args ...interface{}) {
	if h.verbosity >= LevelError {
		fmt.Fprintf(h.err, "[ERROR] "+format+"\n", args...)
	}
}

// FormatError formats an error with masked sensitive data
func (h *Handler) FormatError(err error, output string) string {
	if err == nil {
		return ""
	}

	maskedOutput := h.MaskSensitiveData(output)
	return fmt.Sprintf("Error: %v\nOutput: %s", err, maskedOutput)
}

// FormatSuccess formats successful output with masked sensitive data
func (h *Handler) FormatSuccess(output string) string {
	if output == "" {
		return "Command completed successfully"
	}
	return h.MaskSensitiveData(output)
}

// PrintCommand prints command information based on verbosity level
func (h *Handler) PrintCommand(cmdType string, content map[string]interface{}) {
	if h.verbosity >= LevelDebug {
		h.Debug("Executing command type: %s", cmdType)
		if cmd, ok := content["command"].(string); ok {
			h.Debug("Command: %s", cmd)
		}
		if wd, ok := content["working_dir"].(string); ok {
			h.Debug("Working directory: %s", wd)
		}
	}
}

// PrintResult prints command result based on verbosity level
func (h *Handler) PrintResult(result *interpreter.Result) {
	if result.Success {
		if h.verbosity >= LevelInfo {
			h.Info(h.FormatSuccess(result.Output))
		}
	} else {
		if h.verbosity >= LevelError {
			h.Error(h.FormatError(result.Error, result.Output))
		}
	}
}
