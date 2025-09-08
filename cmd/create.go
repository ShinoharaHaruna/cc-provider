package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Interactively creates a new provider environment.",
	Long:  `Interactively prompts for the necessary details to create a new provider environment file in the ~/.cc-provider directory.`,
	Run:   runCreateCmd,
}

func runCreateCmd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	// 1. Prompt for environment name
	envName := prompt(reader, "Enter environment name (e.g., 'deepseek')", true)
	envFilePath := filepath.Join(cfgDir, envName)

	if _, err := os.Stat(envFilePath); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Environment '%s' already exists.\n", envName)
		os.Exit(1)
	}

	var envVars []string

	// 2. Required inputs
	fmt.Println("\nEnter required variables: ")
	envVars = append(envVars, "ANTHROPIC_BASE_URL="+prompt(reader, "  ANTHROPIC_BASE_URL (e.g., https://api.deepseek.com/anthropic)", true))
	envVars = append(envVars, "ANTHROPIC_AUTH_TOKEN="+prompt(reader, "  ANTHROPIC_AUTH_TOKEN", true))

	// 3. Recommended inputs
	fmt.Println("\nEnter recommended variables (press Enter to skip): ")
	if model := prompt(reader, "  ANTHROPIC_MODEL (e.g., deepseek-chat)", false); model != "" {
		envVars = append(envVars, "ANTHROPIC_MODEL="+model)
	}
	if smallModel := prompt(reader, "  ANTHROPIC_SMALL_FAST_MODEL (e.g., deepseek-coder)", false); smallModel != "" {
		envVars = append(envVars, "ANTHROPIC_SMALL_FAST_MODEL="+smallModel)
	}

	// 4. Optional inputs with defaults
	fmt.Println("\nEnter optional variables (press Enter to use default): ")
	timeout := promptWithDefault(reader, "  API_TIMEOUT_MS", "600000")
	envVars = append(envVars, "API_TIMEOUT_MS="+timeout)

	disableTraffic := promptWithDefault(reader, "  CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1")
	envVars = append(envVars, "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC="+disableTraffic)

	// 5. Write to file
	var finalEnvVars []string
	for _, envVar := range envVars {
		if parts := strings.SplitN(envVar, "=", 2); len(parts) == 2 {
			finalEnvVars = append(finalEnvVars, fmt.Sprintf("%s=\"%s\"", parts[0], parts[1]))
		} else {
			finalEnvVars = append(finalEnvVars, envVar)
		}
	}

	content := strings.Join(finalEnvVars, "\n")
	if err := os.WriteFile(envFilePath, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing environment file '%s': %v\n", envFilePath, err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully created environment '%s'.\n", envName)
	fmt.Printf("To activate it, run: cc-provider activate %s\n", envName)
}

// prompt asks the user for input with a given message.
func prompt(reader *bufio.Reader, message string, required bool) string {
	for {
		fmt.Printf("%s: ", message)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" || !required {
			return input
		}
		fmt.Println("This field is required.")
	}
}

// promptWithDefault asks the user for input with a default value.
func promptWithDefault(reader *bufio.Reader, message, defaultValue string) string {
	fmt.Printf("%s (default: %s): ", message, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func init() {
	rootCmd.AddCommand(createCmd)
}
