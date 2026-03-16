package sandbox

import (
	"fmt"
	"os/exec"
)

// ConfigureProxy configures network proxy to allow the sandbox to reach host ports.
// This is used to allow the sandbox to connect to services like Ollama running on the host.
func (c *SandboxClient) ConfigureProxy(sandboxName string, port int) error {
	cmd := exec.Command("docker", "sandbox", "network", "proxy", sandboxName,
		"--allow-host", fmt.Sprintf("localhost:%d", port))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to configure network proxy: %w, output: %s", err, string(output))
	}

	return nil
}

// ConfigureProxyForOllama configures network proxy for the default Ollama port (11434).
func (c *SandboxClient) ConfigureProxyForOllama(sandboxName string) error {
	return c.ConfigureProxy(sandboxName, 11434)
}
