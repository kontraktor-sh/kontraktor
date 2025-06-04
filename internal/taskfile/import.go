// import.go
// Handles loading, downloading, and cloning of imported taskfiles for Kontraktor.
package taskfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseTaskfile reads and parses a taskfile.ktr.yml from the given path, recursively loading imports.
func ParseTaskfile(path string) (*Taskfile, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("taskfile not found: %s", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open taskfile: %w", err)
	}
	defer f.Close()

	var tf Taskfile
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&tf); err != nil {
		return nil, fmt.Errorf("decode yaml in %s: %w", path, err)
	}

	// Validate the taskfile
	if err := tf.Validate(); err != nil {
		return nil, fmt.Errorf("invalid taskfile: %w", err)
	}

	// Recursively load imports
	for _, importPath := range tf.Imports {
		var importFile string
		if isGitImport(importPath) {
			importFile, err = cloneAndGetFile(importPath)
			if err != nil {
				return nil, fmt.Errorf("git import %s: %w", importPath, err)
			}
		} else if isHTTPImport(importPath) {
			importFile, err = downloadToTemp(importPath)
			if err != nil {
				return nil, fmt.Errorf("download import %s: %w", importPath, err)
			}
		} else {
			importFile = importPath
		}
		imported, err := ParseTaskfile(importFile)
		if err != nil {
			return nil, fmt.Errorf("import %s: %w", importPath, err)
		}
		// Merge imported tasks, but do not override main file tasks
		for k, v := range imported.Tasks {
			if _, exists := tf.Tasks[k]; !exists {
				tf.Tasks[k] = v
			}
		}
	}

	return &tf, nil
}

// isHTTPImport returns true if the path is an HTTP(S) URL.
func isHTTPImport(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// isGitImport returns true if the path is a git repo import (contains .git//).
func isGitImport(path string) bool {
	return strings.Contains(path, ".git//")
}

// downloadToTemp downloads a file from a URL to a temp file and returns its path.
func downloadToTemp(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "ktr-import-*.yml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// cloneAndGetFile clones a git repo and returns the path to the specified file within it.
func cloneAndGetFile(gitImport string) (string, error) {
	parts := strings.SplitN(gitImport, ".git//", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid git import: %s", gitImport)
	}
	repoURL := parts[0] + ".git"
	fileInRepo := parts[1]

	tmpDir, err := os.MkdirTemp("", "ktr-git-*")
	if err != nil {
		return "", err
	}

	fmt.Printf("Cloning repo %s...\n", repoURL)
	cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
	cmd.Stdout = nil // suppress output
	cmd.Stderr = nil // suppress output
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed: %w", err)
	}
	fmt.Printf("Clone complete: %s\n", repoURL)

	fullPath := filepath.Join(tmpDir, fileInRepo)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file %s not found in repo %s", fileInRepo, repoURL)
	}
	return fullPath, nil
}
