package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateRandom(t *testing.T) {
	// Create multiple random workspaces and verify they're unique
	paths := make(map[string]bool)

	for i := 0; i < 10; i++ {
		path, err := CreateRandom()
		if err != nil {
			t.Fatalf("CreateRandom() error: %v", err)
		}

		// Verify path is absolute
		if !filepath.IsAbs(path) {
			t.Errorf("CreateRandom() = %q, want absolute path", path)
		}

		// Verify path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("CreateRandom() created path doesn't exist: %q", path)
		}

		// Verify path is under ~/.cloma/workspaces/
		home, _ := os.UserHomeDir()
		expectedPrefix := filepath.Join(home, ".cloma", "workspaces", "cloma-")
		if len(path) < len(expectedPrefix)+6 {
			t.Errorf("CreateRandom() path too short: %q", path)
		}

		// Verify uniqueness
		if paths[path] {
			t.Errorf("CreateRandom() created duplicate path: %q", path)
		}
		paths[path] = true

		// Clean up
		os.RemoveAll(path)
	}
}

func TestRandomHex(t *testing.T) {
	// Test hex generation
	hex1, err := randomHex(6)
	if err != nil {
		t.Fatalf("randomHex(6) error: %v", err)
	}

	// Should be 12 characters (6 bytes * 2 hex chars per byte)
	if len(hex1) != 12 {
		t.Errorf("randomHex(6) = %q, want 12 characters", hex1)
	}

	// Test uniqueness
	hex2, err := randomHex(6)
	if err != nil {
		t.Fatalf("randomHex(6) error: %v", err)
	}

	if hex1 == hex2 {
		t.Errorf("randomHex not unique: %q == %q", hex1, hex2)
	}

	// Test different lengths
	hex3, err := randomHex(3)
	if err != nil {
		t.Fatalf("randomHex(3) error: %v", err)
	}
	if len(hex3) != 6 {
		t.Errorf("randomHex(3) = %q, want 6 characters", hex3)
	}
}