package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information (set at build time)
var (
	// Version is the semantic version of cloma
	Version = "dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildDate is the date the binary was built
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Print version information for cloma.

This command displays the version, git commit, build date, and Go version
used to build the binary.`,
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// VersionInfo holds version information for output
type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

func runVersion(cmd *cobra.Command, args []string) {
	info := VersionInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	if jsonOutput {
		outputVersionJSON(info)
	} else {
		outputVersionText(info)
	}
}

func outputVersionJSON(info VersionInfo) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(info)
}

func outputVersionText(info VersionInfo) {
	fmt.Printf("cloma version %s\n", info.Version)
	fmt.Printf("  Git commit: %s\n", info.GitCommit)
	fmt.Printf("  Build date: %s\n", info.BuildDate)
	fmt.Printf("  Go version: %s\n", info.GoVersion)
	fmt.Printf("  Platform:   %s\n", info.Platform)
}
