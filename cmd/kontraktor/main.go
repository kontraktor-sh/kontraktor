package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/rafaelherik/kontraktor-sh/internal/taskfile"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kontraktor",
		Short: "Kontraktor CLI",
	}

	var runCmd = &cobra.Command{
		Use:   "run [taskfile] [task]",
		Short: "Run a specific task from a taskfile.ktr.yml",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			taskName := args[1]
			tf, err := taskfile.ParseTaskfile(path)
			if err != nil {
				return fmt.Errorf("failed to parse taskfile: %w", err)
			}
			task, ok := tf.Tasks[taskName]
			if !ok {
				return fmt.Errorf("task '%s' not found", taskName)
			}
			fmt.Printf("Running task '%s': %s\n", taskName, task.Desc)
			for i, cmdStr := range task.Cmds {
				fmt.Printf("[%d/%d] $ %s\n", i+1, len(task.Cmds), cmdStr)
				if err := runShellCommand(cmdStr); err != nil {
					return fmt.Errorf("command failed: %w", err)
				}
			}
			return nil
		},
	}

	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runShellCommand(cmdStr string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
} 