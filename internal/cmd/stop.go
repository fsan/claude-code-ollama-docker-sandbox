package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"cloma/internal/config"
	"cloma/internal/sandbox"
	"cloma/internal/workspace"
)

var stopWorkspace string

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running sandbox",
	Long: `Stop a running Docker sandbox.

This command stops the sandbox for the specified workspace.
The sandbox is preserved and can be restarted later with 'cloma run' or 'cloma shell'.

If the sandbox does not exist or is already stopped, this command does nothing.`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().StringVarP(&stopWorkspace, "workspace", "w", "", "Workspace directory (default: current directory)")
}

func runStop(cmd *cobra.Command, args []string) error {
	// Initialize config
	if err := config.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Resolve workspace
	workspacePath := stopWorkspace
	if workspacePath == "" {
		workspacePath = "."
	}

	resolvedWorkspace, err := workspace.Resolve(workspacePath)
	if err != nil {
		return fmt.Errorf("failed to resolve workspace: %w\nHint: Ensure the path exists: %s", err, workspacePath)
	}

	// Generate sandbox name
	sandboxName := workspace.SandboxName(resolvedWorkspace)

	// Check prerequisites
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed or not in PATH\nHint: Install Docker Desktop from https://www.docker.com/products/docker-desktop")
	}

	if err := sandbox.EnsureSandboxPlugin(); err != nil {
		return fmt.Errorf("Docker Desktop sandbox plugin required\nHint: Requires Docker Desktop 4.58+\nEnable sandbox plugin in Docker Desktop settings")
	}

	// Check if sandbox exists
	exists, err := sandbox.Exists(sandboxName)
	if err != nil {
		return fmt.Errorf("failed to check if sandbox exists: %w", err)
	}

	if !exists {
		fmt.Printf("Sandbox does not exist: %s\n", sandboxName)
		return nil
	}

	// Check if running
	isRunning, err := sandbox.IsRunning(sandboxName)
	if err != nil {
		return fmt.Errorf("failed to check sandbox status: %w", err)
	}

	if !isRunning {
		fmt.Printf("Sandbox is not running: %s\n", sandboxName)
		return nil
	}

	// Stop the sandbox
	fmt.Printf("Stopping sandbox: %s\n", sandboxName)
	if err := sandbox.Stop(sandboxName); err != nil {
		return fmt.Errorf("failed to stop sandbox: %w", err)
	}

	fmt.Println("Sandbox stopped.")
	return nil
}