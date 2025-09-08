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

	// 2. Ensure shell config is set up
	shellRCPath, err := ensureShellConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up shell configuration: %v\n", err)
		os.Exit(1)
	}

	// 3. Generate and write active_env.sh
	if err := writeActiveEnvScript(envName, envFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing active environment script: %v\n", err)
		os.Exit(1)
	}

	// 4. Prompt user
	fmt.Printf("Successfully activated environment '%s'.\n", envName)
	fmt.Println("Please run the following command to apply changes to your current shell session, or open a new terminal:")
	fmt.Printf("  source %s\n", shellRCPath)
}

// ensureShellConfig checks and modifies the user's shell configuration file.
func ensureShellConfig() (string, error) {
	shell := os.Getenv("SHELL")
	var rcFileName string

	if strings.Contains(shell, "zsh") {
		rcFileName = ".zshrc"
	} else if strings.Contains(shell, "bash") {
		rcFileName = ".bashrc"
	} else {
		fmt.Println("Could not detect shell. Please manually add the following line to your shell's config file:")
		fmt.Printf("  source \"%s\"\n", activeEnvFile)
		return "your_shell_config_file", nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	rcPath := filepath.Join(home, rcFileName)

	// Check if the source line already exists
	file, err := os.Open(rcPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("opening %s: %w", rcPath, err)
		}
		// If file does not exist, we will create it later.
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		sourceLine := fmt.Sprintf("source \"%s\"", activeEnvFile)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), sourceLine) {
				return rcPath, nil // Already configured
			}
		}
	}

	// Append the source line to the rc file
	f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("opening %s for appending: %w", rcPath, err)
	}
	defer f.Close()

	sourceLine := fmt.Sprintf(`
# Added by cc-provider to source active environment variables
source "%s"
`, activeEnvFile)

	if _, err := f.WriteString(sourceLine); err != nil {
		return "", fmt.Errorf("writing to %s: %w", rcPath, err)
	}

	fmt.Printf("Added source command to '%s' to automatically load environments.\n", rcPath)
	return rcPath, nil
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
