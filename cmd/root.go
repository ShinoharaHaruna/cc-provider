package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// cfgDir is the configuration directory for cc-provider.
	cfgDir string

	// activeEnvFile is the path to the script that holds the active environment variables.
	activeEnvFile string

	// envVarKeys holds all the environment variable keys that cc-provider manages.
	envVarKeys = []string{
		"ANTHROPIC_BASE_URL",
		"ANTHROPIC_AUTH_TOKEN",
		"ANTHROPIC_MODEL",
		"ANTHROPIC_SMALL_FAST_MODEL",
		"API_TIMEOUT_MS",
		"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC",
		"CC_PROVIDER_ACTIVE_ENV",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cc-provider",
	Short: "A CLI tool to manage Claude Code provider environments.",
	Long: `cc-provider is a command-line interface tool designed to help you manage
different sets of environment variables for Claude Code. You can create, list,
activate, and remove environments with ease.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig creates the configuration directory if it does not exist.
func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	cfgDir = filepath.Join(home, ".cc-provider")
	activeEnvFile = filepath.Join(cfgDir, "active_env.sh")

	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory '%s': %v\n", cfgDir, err)
			os.Exit(1)
		}
	}
}
