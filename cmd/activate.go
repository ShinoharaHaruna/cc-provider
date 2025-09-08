package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var activateCmd = &cobra.Command{
	Use:   "activate [env-name]",
	Short: "Activates a provider environment.",
	Long: `Activates a specified provider environment by updating the active environment script.
This command also ensures that the user's shell configuration file sources the script.`,
	Args: cobra.ExactArgs(1),
	Run:  runActivateCmd,
}

func runActivateCmd(cmd *cobra.Command, args []string) {
	envName := args[0]
	envFilePath := filepath.Join(cfgDir, envName)

	// 1. Validate environment exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Environment '%s' not found.\n", envName)
		os.Exit(1)
	}

	// 2. Generate and write active_env.sh
	if err := writeActiveEnvScript(envName, envFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing active environment script: %v\n", err)
		os.Exit(1)
	}

	// 3. Prompt user
	fmt.Printf("Successfully activated environment '%s'.\n", envName)
	fmt.Println("Please restart your shell or source your shell's config file to apply the changes.")
}

// writeActiveEnvScript generates the content for and writes to the active_env.sh file.
func writeActiveEnvScript(envName, envFilePath string) error {
	var sb strings.Builder

	// Always start by unsetting all managed keys to ensure a clean state.
	sb.WriteString("# Unset previous variables managed by cc-provider\n")
	for _, key := range envVarKeys {
		sb.WriteString(fmt.Sprintf("unset %s\n", key))
	}
	sb.WriteString("\n")

	// Read variables from the environment file
	content, err := os.ReadFile(envFilePath)
	if err != nil {
		return fmt.Errorf("reading env file %s: %w", envFilePath, err)
	}

	// Add export commands for the new environment
	sb.WriteString(fmt.Sprintf("# Export variables for environment: %s\n", envName))
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			sb.WriteString(fmt.Sprintf("export %s=%s\n", key, value))
		}
	}
	sb.WriteString("\n")

	// Set the active environment identifier
	sb.WriteString(fmt.Sprintf(`export CC_PROVIDER_ACTIVE_ENV="%s"`, envName))
	sb.WriteString("\n")

	// Write the generated script to the active_env.sh file.
	return os.WriteFile(activeEnvFile, []byte(sb.String()), 0644)
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
