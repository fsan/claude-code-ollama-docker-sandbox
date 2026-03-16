package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fsan/cloma/internal/sandbox"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cloma-managed sandboxes",
	Long: `List all Docker Desktop sandboxes managed by cloma.

Sandboxes managed by cloma have names starting with "cloma-".
The list shows the sandbox name, status, and the decoded workspace path.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// SandboxInfo holds information about a sandbox for display
type SandboxInfo struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Workspace string `json:"workspace,omitempty"`
}

func runList(cmd *cobra.Command, args []string) error {
	// Get all sandboxes
	sandboxes, err := sandbox.List()
	if err != nil {
		return fmt.Errorf("failed to list sandboxes: %w", err)
	}

	// Filter to cloma-managed sandboxes (names starting with "cloma-")
	var clomaSandboxes []SandboxInfo
	for _, sb := range sandboxes {
		if strings.HasPrefix(sb.Name, "cloma-") {
			info := SandboxInfo{
				Name:   sb.Name,
				Status: sb.Status,
			}
			// Try to decode workspace from name
			info.Workspace = decodeWorkspaceFromName(sb.Name)
			clomaSandboxes = append(clomaSandboxes, info)
		}
	}

	// Output
	if jsonOutput {
		return outputJSON(clomaSandboxes)
	}

	return outputText(clomaSandboxes)
}

func outputJSON(sandboxes []SandboxInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(sandboxes)
}

func outputText(sandboxes []SandboxInfo) error {
	if len(sandboxes) == 0 {
		fmt.Println("No cloma-managed sandboxes found.")
		return nil
	}

	// Print header
	fmt.Printf("%-50s %-12s %s\n", "NAME", "STATUS", "WORKSPACE")
	fmt.Println(strings.Repeat("-", 80))

	// Print sandboxes
	for _, sb := range sandboxes {
		workspace := sb.Workspace
		if workspace == "" {
			workspace = "<unknown>"
		}
		fmt.Printf("%-50s %-12s %s\n", sb.Name, sb.Status, workspace)
	}

	return nil
}

// decodeWorkspaceFromName attempts to decode the workspace path from a sandbox name.
// Sandbox names follow the pattern: cloma-{slug}-{hash}
// The slug is derived from the workspace path basename.
func decodeWorkspaceFromName(name string) string {
	// Remove "cloma-" prefix
	if !strings.HasPrefix(name, "cloma-") {
		return ""
	}

	// Extract the slug part (everything after prefix, minus the hash suffix)
	parts := strings.TrimPrefix(name, "cloma-")

	// The hash is 8 characters at the end
	if len(parts) < 10 { // minimum: slug + "-" + 8-char hash
		return ""
	}

	// Find the last hyphen to separate slug from hash
	lastHyphen := strings.LastIndex(parts, "-")
	if lastHyphen == -1 {
		return ""
	}

	slug := parts[:lastHyphen]
	return slug
}
