package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Re-run the shell integration setup.",
	Long: `Re-run the shell integration setup to ensure shell functions and completion are properly installed.
This is useful after upgrading cc-provider or if the shell integration is not working.`,
	Run: runSetupCmd,
}

func runSetupCmd(cmd *cobra.Command, args []string) {
	// Force re-run the shell configuration setup
	// 强制重新运行 shell 配置设置
	rcPath, err := ensureShellConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during shell setup: %v\n", err)
		os.Exit(1)
	}

	if rcPath != "" {
		fmt.Printf("Shell integration setup complete.\n")
		fmt.Printf("Configuration updated in '%s'.\n", rcPath)
		fmt.Println("\nPlease restart your shell or run:")
		fmt.Printf("  source %s\n", rcPath)
	} else {
		fmt.Println("Shell integration setup complete.")
		fmt.Println("Your shell is not automatically supported, but you can manually source:")
		fmt.Printf("  source %s\n", activeEnvFile)
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
