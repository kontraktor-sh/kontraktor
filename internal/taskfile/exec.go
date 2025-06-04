// Package taskfile provides task execution logic for Kontraktor.
package taskfile

import "fmt"

// ExecuteTask executes a task by name, tracking visited tasks to detect cycles.
func ExecuteTask(taskName string, tf *Taskfile, visited map[string]bool) error {
	if visited[taskName] {
		return fmt.Errorf("circular reference detected at task '%s'", taskName)
	}
	visited[taskName] = true
	task, ok := tf.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task '%s' not found", taskName)
	}
	for _, cmd := range task.Cmds {
		if cmd.Type == "task" {
			if name, ok := cmd.Content["name"].(string); ok && name != "" {
				if err := ExecuteTask(name, tf, visited); err != nil {
					return err
				}
			}
		}
	}
	visited[taskName] = false
	return nil
}
