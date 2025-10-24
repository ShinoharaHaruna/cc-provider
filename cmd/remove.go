package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:               "remove [env-name]",
	Short:             "Removes a provider environment.",
	Long:              `Removes a specified provider environment file. If the environment is currently active, it will be deactivated.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeEnvironmentNamesForRemove,
	Run:               runRemoveCmd,
}

func runRemoveCmd(cmd *cobra.Command, args []string) {
	envName := args[0]
	envFilePath := filepath.Join(cfgDir, envName)

	// 1. Validate environment exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Environment '%s' not found. // 错误: 未找到环境 '%s'。\n", envName, envName)
		os.Exit(1)
	}

	// 2. Remove the environment file
	if err := os.Remove(envFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing environment file '%s': %v\n", envFilePath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully removed environment '%s'.\n", envName)

	// 3. Check if the removed environment was the active one
	activeEnv := os.Getenv("CC_PROVIDER_ACTIVE_ENV")
	if activeEnv == envName {
		if err := deactivateActiveEnv(); err != nil {
			fmt.Fprintf(os.Stderr, "Error deactivating environment: %v\n", err)
			os.Exit(1)
		}

		shellRCPath := "your_shell_config_file"
		shell := os.Getenv("SHELL")
		if strings.Contains(shell, "zsh") {
			shellRCPath = "~/.zshrc"
		} else if strings.Contains(shell, "bash") {
			shellRCPath = "~/.bashrc"
		}

		fmt.Printf("Deactivated environment '%s'.\n", envName)
		fmt.Println("Please run the following command to apply changes, or open a new terminal:")
		fmt.Printf("  source %s\n", shellRCPath)
	}
}

// deactivateActiveEnv clears the active_env.sh script.
func deactivateActiveEnv() error {
	var sb strings.Builder

	// Add unset commands for all managed keys
	sb.WriteString("# Unset previous variables managed by cc-provider\n")
	for _, key := range envVarKeys {
		sb.WriteString(fmt.Sprintf("unset %s\n", key))
	}
	sb.WriteString("\n")

	// Set the active environment identifier to empty
	sb.WriteString(`export CC_PROVIDER_ACTIVE_ENV=""`)
	sb.WriteString("\n")

	// Write the script to file
	return os.WriteFile(activeEnvFile, []byte(sb.String()), 0644)
}

// completeEnvironmentNamesForRemove provides completion for environment names
// 为环境名称提供补全
func completeEnvironmentNamesForRemove(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return getEnvironmentNames(), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
