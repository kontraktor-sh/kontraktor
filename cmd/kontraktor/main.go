package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kontraktor-sh/kontraktor/internal/taskfile"
	"github.com/kontraktor-sh/kontraktor/internal/vault"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kontraktor",
		Short: "Kontraktor CLI",
	}

	var fileFlag string
	var runCmd = &cobra.Command{
		Use:   "run [task] [key=value]...",
		Short: "Run a specific task with arguments from a .ktr.yml file in the current directory or specified with -f",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]
			argMap := make(map[string]string)
			for _, arg := range args[1:] {
				kv := strings.SplitN(arg, "=", 2)
				if len(kv) != 2 {
					return fmt.Errorf("invalid argument: %s (expected key=value)", arg)
				}
				argMap[kv[0]] = kv[1]
			}
			var path string
			if fileFlag != "" {
				path = fileFlag
			} else {
				files, err := findKtrYAMLFiles()
				if err != nil {
					return fmt.Errorf("failed to search for .ktr.yml files: %w", err)
				}
				if len(files) == 0 {
					return fmt.Errorf("no .ktr.yml files found in current directory")
				}
				if len(files) > 1 {
					fmt.Println("Multiple .ktr.yml files found:")
					for _, f := range files {
						fmt.Println(" -", f)
					}
					return fmt.Errorf("please specify the file explicitly with -f")
				}
				path = files[0]
			}
			tf, err := taskfile.ParseTaskfile(path)
			if err != nil {
				return fmt.Errorf("failed to parse taskfile: %w", err)
			}
			visited := make(map[string]bool)
			if err := executeTaskWithArgs(taskName, tf, visited, argMap, nil, []int{}, nil); err != nil {
				return err
			}
			return nil
		},
	}

	runCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to .ktr.yml file")

	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runShellCommand(cmdStr string, env map[string]string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stderr = os.Stderr
	if env != nil {
		cmd.Env = mergeWithOSEnv(env)
	}

	// Create pipes for stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Read and mask stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Mask any secret values in the output
		for _, v := range env {
			if strings.Contains(line, v) {
				line = strings.ReplaceAll(line, v, "<<SENSITIVE>>")
			}
		}
		fmt.Println(line)
	}

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return scanner.Err()
}

// Helper to merge custom env with os.Environ
func mergeWithOSEnv(env map[string]string) []string {
	base := os.Environ()
	result := make([]string, 0, len(base)+len(env))
	seen := make(map[string]struct{})
	for _, e := range base {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) == 2 {
			if _, ok := env[kv[0]]; ok {
				continue // overridden
			}
			seen[kv[0]] = struct{}{}
			result = append(result, e)
		}
	}
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}

func findKtrYAMLFiles() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".ktr.yml") {
			files = append(files, name)
		}
	}
	return files, nil
}

func executeTaskWithArgs(taskName string, tf *taskfile.Taskfile, visited map[string]bool, args map[string]string, parentArgs map[string]string, stepPath []int, parentEnv map[string]string) error {
	if visited[taskName] {
		return fmt.Errorf("circular reference detected at task '%s'", taskName)
	}
	visited[taskName] = true
	task, ok := tf.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task '%s' not found", taskName)
	}
	// Merge parentArgs and args, with args taking precedence
	mergedArgs := make(map[string]string)
	for k, v := range parentArgs {
		mergedArgs[k] = v
	}
	for k, v := range args {
		mergedArgs[k] = v
	}
	// Fill in defaults for missing args
	for _, arg := range task.Args {
		if _, ok := mergedArgs[arg.Name]; !ok {
			if arg.Default != nil {
				mergedArgs[arg.Name] = fmt.Sprintf("%v", arg.Default)
			} else {
				return fmt.Errorf("missing required argument: %s", arg.Name)
			}
		}
	}
	// Merge environment: global -> parentEnv -> per-task
	mergedEnv := make(map[string]string)
	for k, v := range tf.Environment {
		mergedEnv[k] = v
	}
	for k, v := range parentEnv {
		mergedEnv[k] = v
	}
	for k, v := range task.Environment {
		mergedEnv[k] = v
	}
	// Fetch Azure Key Vault secrets if configured
	ctx := context.Background()
	if tf.Vaults != nil && tf.Vaults.AzureKeyVault != nil {
		for _, config := range tf.Vaults.AzureKeyVault {
			secrets, err := vault.FetchAzureSecrets(ctx, config)
			if err != nil {
				return fmt.Errorf("failed to fetch Azure secrets: %w", err)
			}
			for k, v := range secrets {
				mergedEnv[k] = v
			}
		}
	}
	fmt.Printf("Running task '%s': %s\n", taskName, task.Desc)
	for i, cmd := range task.Cmds {
		currentStepPath := append(stepPath, i+1)
		stepNum := formatStepNum(currentStepPath, len(task.Cmds))
		if cmd.Task != "" {
			fmt.Printf("[%s] > task: %s\n", stepNum, cmd.Task)
			if err := executeTaskWithArgs(cmd.Task, tf, visited, mergedArgs, mergedArgs, append(stepPath, i+1), mergedEnv); err != nil {
				return err
			}
			continue
		}
		if cmd.Uses != "" {
			fmt.Printf("[%s] > uses: %s\n", stepNum, cmd.Uses)
			usesArgs := make(map[string]string)
			for k, v := range mergedArgs {
				usesArgs[k] = v
			}
			for k, v := range cmd.Args {
				usesArgs[k] = fmt.Sprintf("%v", v)
			}
			if err := executeTaskWithArgs(cmd.Uses, tf, visited, usesArgs, usesArgs, append(stepPath, i+1), mergedEnv); err != nil {
				return err
			}
			continue
		}
		cmdStr := cmd.Cmd
		for k, v := range mergedArgs {
			cmdStr = strings.ReplaceAll(cmdStr, "${"+k+"}", v)
		}
		// Handle secrets substitution
		if tf.Vaults != nil && tf.Vaults.AzureKeyVault != nil {
			for _, config := range tf.Vaults.AzureKeyVault {
				secrets, err := vault.FetchAzureSecrets(ctx, config)
				if err != nil {
					return fmt.Errorf("failed to fetch Azure secrets: %w", err)
				}
				for k, v := range secrets {
					cmdStr = strings.ReplaceAll(cmdStr, "${secrets."+k+"}", v)
				}
			}
		}
		fmt.Printf("[%s] $ %s\n", stepNum, cmdStr)
		if err := runShellCommand(cmdStr, mergedEnv); err != nil {
			return fmt.Errorf("command failed: %w", err)
		}
	}
	visited[taskName] = false
	return nil
}

// Helper to format the step number as e.g. 2/3, 2.1/2, 3/3, 3.1/1
func formatStepNum(stepPath []int, total int) string {
	if len(stepPath) == 0 {
		return ""
	}
	var parts []string
	for _, n := range stepPath {
		parts = append(parts, fmt.Sprintf("%d", n))
	}
	return fmt.Sprintf("%s/%d", strings.Join(parts, "."), total)
}
