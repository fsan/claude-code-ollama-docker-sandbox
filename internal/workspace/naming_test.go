package workspace

import (
	"testing"
)

func TestPathToSlug(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/home/user/myproject",
			expected: "myproject",
		},
		{
			name:     "path with spaces",
			path:     "/home/user/My Project",
			expected: "my-project",
		},
		{
			name:     "path with special chars",
			path:     "/home/user/my-project_2024",
			expected: "my-project-2024",
		},
		{
			name:     "path with mixed case",
			path:     "/home/user/MyProject",
			expected: "myproject",
		},
		{
			name:     "path with trailing special char",
			path:     "/home/user/project!",
			expected: "project",
		},
		{
			name:     "path with leading special char",
			path:     "/home/user/.hidden",
			expected: "hidden",
		},
		{
			name:     "path with multiple special chars",
			path:     "/home/user/My___Project!!!",
			expected: "my-project",
		},
		{
			name:     "current directory",
			path:     ".",
			expected: "", // "." becomes "-", which gets stripped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PathToSlug(tt.path)
			if result != tt.expected {
				t.Errorf("PathToSlug(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestPathHash(t *testing.T) {
	// Test that hashes are consistent
	path1 := "/home/user/project"
	hash1 := PathHash(path1)
	hash2 := PathHash(path1)

	if hash1 != hash2 {
		t.Errorf("PathHash not consistent: %q != %q", hash1, hash2)
	}

	// Test that hashes are 8 characters
	if len(hash1) != 8 {
		t.Errorf("PathHash length = %d, want 8", len(hash1))
	}

	// Test that different paths produce different hashes
	path2 := "/home/user/other"
	hash3 := PathHash(path2)
	if hash1 == hash3 {
		t.Errorf("PathHash collision: %q and %q both produce %q", path1, path2, hash1)
	}
}

func TestSandboxName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		contains string
	}{
		{
			name:     "simple path",
			path:     "/home/user/myproject",
			contains: "cloma-myproject-",
		},
		{
			name:     "path with spaces",
			path:     "/home/user/My Project",
			contains: "cloma-my-project-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SandboxName(tt.path)
			// Verify length is adequate (cloma- + slug + - + 8-char hash)
			if len(result) < 15 {
				t.Errorf("SandboxName(%q) = %q, too short", tt.path, result)
			}
			// Check prefix
			prefix := "cloma-"
			if result[:len(prefix)] != prefix {
				t.Errorf("SandboxName(%q) = %q, should start with %q", tt.path, result, prefix)
			}
			// Verify it contains the expected substring
			if len(result) >= len(tt.contains) && result[:len(tt.contains)] != tt.contains {
				t.Errorf("SandboxName(%q) = %q, should start with %q", tt.path, result, tt.contains)
			}
		})
	}
}
