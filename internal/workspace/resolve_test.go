package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "workspace-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		setup    func() error
		cleanup  func() error
		wantErr  bool
		check    func(t *testing.T, result string)
	}{
		{
			name:    "empty string creates random workspace",
			input:   "",
			wantErr: false,
			check: func(t *testing.T, result string) {
				// Should be in ~/.cloma/workspaces/
				if !filepath.IsAbs(result) {
					t.Errorf("Result is not absolute: %q", result)
				}
				expectedPrefix := filepath.Join(home, ".cloma", "workspaces", "cloma-")
				if len(result) < len(expectedPrefix)+6 {
					t.Errorf("Result too short: %q", result)
				}
			},
		},
		{
			name:    "tilde resolves to home",
			input:   "~",
			wantErr: false,
			check: func(t *testing.T, result string) {
				if result != home {
					t.Errorf("Resolve(~) = %q, want %q", result, home)
				}
			},
		},
		{
			name:    "tilde with subpath",
			input:   "~/subdir",
			wantErr: true, // The path doesn't exist
		},
		{
			name:    "dot resolves to cwd",
			input:   ".",
			wantErr: false,
			check: func(t *testing.T, result string) {
				if result != cwd {
					t.Errorf("Resolve(.) = %q, want %q", result, cwd)
				}
			},
		},
		{
			name:    "absolute path that exists",
			input:   tmpDir,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if result != tmpDir {
					t.Errorf("Resolve(%q) = %q, want %q", tmpDir, result, tmpDir)
				}
			},
		},
		{
			name:    "relative path that exists",
			input:   ".", // current directory (workspace package)
			wantErr: false,
			check: func(t *testing.T, result string) {
				// Should resolve to the workspace package directory
				if !filepath.IsAbs(result) {
					t.Errorf("Result is not absolute: %q", result)
				}
			},
		},
		{
			name:    "non-existent path",
			input:   "/non/existent/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			result, err := Resolve(tt.input)

			if tt.cleanup != nil {
				if err := tt.cleanup(); err != nil {
					t.Fatalf("Cleanup failed: %v", err)
				}
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("Resolve(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Resolve(%q) unexpected error: %v", tt.input, err)
				} else if tt.check != nil {
					tt.check(t, result)
				}
			}

			// Clean up random workspace if created
			if tt.input == "" && result != "" {
				os.RemoveAll(result)
			}
		})
	}
}