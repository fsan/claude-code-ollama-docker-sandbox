package sandbox

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Sandbox represents a Docker Desktop sandbox instance.
type Sandbox struct {
	// Name is the unique identifier for the sandbox.
	Name string `json:"name"`

	// Agent is the agent type for the sandbox.
	Agent string `json:"agent"`

	// Status indicates the current state of the sandbox (e.g., "running", "stopped").
	Status string `json:"status"`

	// Image is the base image used for the sandbox.
	Image string `json:"image"`
}

// sandboxListResponse represents the JSON response from `docker sandbox ls --json`.
type sandboxListResponse struct {
	VMs []Sandbox `json:"vms"`
}

// List returns all Docker Desktop sandboxes.
// It parses the JSON output from `docker sandbox ls --json`.
func List() ([]Sandbox, error) {
	cmd := exec.Command("docker", "sandbox", "ls", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list sandboxes: %w", err)
	}

	var response sandboxListResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse sandbox list: %w", err)
	}

	return response.VMs, nil
}

// Exists checks if a sandbox with the given name exists.
func Exists(name string) (bool, error) {
	sandboxes, err := List()
	if err != nil {
		return false, err
	}

	for _, sb := range sandboxes {
		if sb.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// IsRunning checks if a sandbox with the given name is currently running.
func IsRunning(name string) (bool, error) {
	sandboxes, err := List()
	if err != nil {
		return false, err
	}

	for _, sb := range sandboxes {
		if sb.Name == name {
			return sb.Status == "running", nil
		}
	}

	return false, nil
}

// Get retrieves a specific sandbox by name.
// Returns nil if the sandbox doesn't exist.
func Get(name string) (*Sandbox, error) {
	sandboxes, err := List()
	if err != nil {
		return nil, err
	}

	for _, sb := range sandboxes {
		if sb.Name == name {
			return &sb, nil
		}
	}

	return nil, nil
}
