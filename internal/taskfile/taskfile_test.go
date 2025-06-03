package taskfile

import (
	"os"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "taskfile-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmp.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()
	t.Cleanup(func() { os.Remove(tmp.Name()) })
	return tmp.Name()
}

func TestParseTaskfile_Scenarios(t *testing.T) {
	t.Run("basic parse", func(t *testing.T) {
		content := `
version: 0.3
tasks:
  hello:
    desc: "Say hello"
    cmds:
      - echo "Hello"
`
		file := writeTempFile(t, content)
		tf, err := ParseTaskfile(file)
		if err != nil {
			t.Fatalf("ParseTaskfile failed: %v", err)
		}
		if tf.Version != "0.3" {
			t.Errorf("expected version 0.3, got %s", tf.Version)
		}
		if _, ok := tf.Tasks["hello"]; !ok {
			t.Errorf("expected task 'hello' to be present")
		}
	})

	t.Run("local import", func(t *testing.T) {
		imported := `
version: 0.3
tasks:
  imported:
    desc: "Imported task"
    cmds:
      - echo "Imported"
`
		importFile := writeTempFile(t, imported)
		main := `
version: 0.3
imports:
  - ` + importFile + `
tasks:
  main:
    desc: "Main task"
    cmds:
      - echo "Main"
`
		mainFile := writeTempFile(t, main)
		tf, err := ParseTaskfile(mainFile)
		if err != nil {
			t.Fatalf("ParseTaskfile with import failed: %v", err)
		}
		if _, ok := tf.Tasks["imported"]; !ok {
			t.Errorf("expected imported task to be present")
		}
	})

	t.Run("circular reference", func(t *testing.T) {
		content := `
version: 0.3
tasks:
  a:
    desc: "A"
    cmds:
      - task: b
  b:
    desc: "B"
    cmds:
      - task: a
`
		file := writeTempFile(t, content)
		tf, err := ParseTaskfile(file)
		if err != nil {
			t.Fatalf("ParseTaskfile failed: %v", err)
		}
		visited := make(map[string]bool)
		err = executeTask("a", tf, visited)
		if err == nil || err.Error() != "circular reference detected at task 'a'" {
			t.Errorf("expected circular reference error, got: %v", err)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := ParseTaskfile("/tmp/does-not-exist-xyz.yml")
		if err == nil || !os.IsNotExist(err) && !contains(err.Error(), "not found") {
			t.Errorf("expected not found error, got: %v", err)
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		file := writeTempFile(t, "not: [valid: yaml")
		_, err := ParseTaskfile(file)
		if err == nil || !contains(err.Error(), "decode yaml") {
			t.Errorf("expected decode yaml error, got: %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
} 