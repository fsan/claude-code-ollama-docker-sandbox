// Package workspace provides workspace management functionality for cloma.
// It handles path resolution, naming conventions, and random workspace creation.
package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"regexp"
	"strings"
)

// PathToSlug converts a path basename to a slug format.
// It lowercases the basename and replaces non-alphanumeric characters with hyphens.
// Consecutive hyphens are collapsed, and leading/trailing hyphens are removed.
func PathToSlug(path string) string {
	basename := filepath.Base(path)

	// Convert to lowercase
	slug := strings.ToLower(basename)

	// Replace non-alphanumeric characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9]`)
	slug = re.ReplaceAllString(slug, "-")

	// Collapse consecutive hyphens
	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")

	// Remove leading hyphen
	slug = strings.TrimPrefix(slug, "-")

	// Remove trailing hyphen
	slug = strings.TrimSuffix(slug, "-")

	return slug
}

// PathHash generates an 8-character SHA256 hash of the given path.
// This provides uniqueness for sandbox names when combined with the slug.
func PathHash(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:])[:8]
}

// SandboxName generates a sandbox name from a workspace path.
// The format is: cloma-{slug}-{hash}
// Example: cloma-myproject-a1b2c3d4
func SandboxName(workspace string) string {
	slug := PathToSlug(workspace)
	hash := PathHash(workspace)
	return "cloma-" + slug + "-" + hash
}
