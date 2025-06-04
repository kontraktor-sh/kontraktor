package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kontraktor-sh/kontraktor/internal/output"
)

// VerbosityLevel represents the output verbosity level
type VerbosityLevel string

const (
	// VerbositySilent suppresses all output except critical errors
	VerbositySilent VerbosityLevel = "SILENT"
	// VerbosityError shows only errors
	VerbosityError VerbosityLevel = "ERROR"
	// VerbosityInfo shows normal output and errors
	VerbosityInfo VerbosityLevel = "INFO"
	// VerbosityDebug shows all output including debug information
	VerbosityDebug VerbosityLevel = "DEBUG"
)

// Config holds the CLI configuration
type Config struct {
	Verbosity    VerbosityLevel
	TaskName     string
	TaskArgs     map[string]string
	MaskPatterns []string
}

// ParseFlags parses command line flags and returns the configuration
func ParseFlags() (*Config, error) {
	config := &Config{
		TaskArgs:     make(map[string]string),
		MaskPatterns: []string{},
	}

	// Parse verbosity flag
	verbosity := flag.String("verbosity", string(VerbosityInfo), "Output verbosity level (SILENT, ERROR, INFO, DEBUG)")
	flag.Parse()

	// Set verbosity level
	switch strings.ToUpper(*verbosity) {
	case string(VerbositySilent):
		config.Verbosity = VerbositySilent
	case string(VerbosityError):
		config.Verbosity = VerbosityError
	case string(VerbosityInfo):
		config.Verbosity = VerbosityInfo
	case string(VerbosityDebug):
		config.Verbosity = VerbosityDebug
	default:
		return nil, fmt.Errorf("invalid verbosity level: %s", *verbosity)
	}

	// Get task name and arguments
	args := flag.Args()
	if len(args) < 2 || args[0] != "run" {
		return nil, fmt.Errorf("usage: kontraktor run <taskname> [args...]")
	}

	config.TaskName = args[1]

	// Parse task arguments (key=value pairs)
	for _, arg := range args[2:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid argument format: %s (expected key=value)", arg)
		}
		config.TaskArgs[parts[0]] = parts[1]
	}

	// Add default mask patterns for sensitive data
	config.MaskPatterns = append(config.MaskPatterns,
		"password=.*",
		"secret=.*",
		"key=.*",
		"token=.*",
	)

	return config, nil
}

// CreateOutputHandler creates an output handler based on the configuration
func (c *Config) CreateOutputHandler() (*output.Handler, error) {
	handler := output.NewHandler()

	// Set verbosity level
	switch c.Verbosity {
	case VerbositySilent:
		handler.SetLevel(output.LevelSilent)
	case VerbosityError:
		handler.SetLevel(output.LevelError)
	case VerbosityInfo:
		handler.SetLevel(output.LevelInfo)
	case VerbosityDebug:
		handler.SetLevel(output.LevelDebug)
	}

	// Add mask patterns
	for _, pattern := range c.MaskPatterns {
		handler.AddMaskPattern(pattern)
	}

	return handler, nil
}
