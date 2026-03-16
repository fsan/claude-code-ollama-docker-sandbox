package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrWorkspaceNotFound is returned when a specified workspace path does not exist.
var ErrWorkspaceNotFound = errors.New("workspace path does not exist")

// Resolve resolves a workspace path string to an absolute path.
// It handles the following cases:
//   - "." or empty string: resolves to the current working directory
//   - "~": resolves to the user's home directory
//   - named paths: resolves relative paths to absolute paths
//   - absolute paths: returned as-is after verification
//
// If no workspace is specified (empty string), a random workspace is created.
func Resolve(workspace string) (string, error) {
	// If no workspace specified, create a random one
	if workspace == "" {
		return CreateRandom()
	}

	// Handle home directory shortcut
	if workspace == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return resolvePath(home)
	}

	// Handle paths starting with ~/
	if strings.HasPrefix(workspace, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		expanded := filepath.Join(home, workspace[2:])
		return resolvePath(expanded)
	}

	// Handle current directory
	if workspace == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return resolvePath(cwd)
	}

	// Handle relative and absolute paths
	return resolvePath(workspace)
}

// resolvePath converts a path to an absolute path and verifies it exists.
func resolvePath(path string) (string, error) {
	// Convert to absolute path if not already
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Verify the path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", ErrWorkspaceNotFound
	}

	return absPath, nil
}
