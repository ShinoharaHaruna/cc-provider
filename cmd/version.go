package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the semantic version of cc-provider
	// 语义化版本号
	Version = "0.1.0"

	// BuildTime is the time when the binary was built
	// 构建时间
	BuildTime = "unknown"

	// GitCommit is the git commit hash
	// Git 提交哈希
	GitCommit = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information.",
	Long:  `Display the version, build time, and git commit information of cc-provider.`,
	Run:   runVersionCmd,
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("cc-provider version %s\n", Version)
	fmt.Printf("Build time: %s\n", BuildTime)
	fmt.Printf("Git commit: %s\n", GitCommit)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
