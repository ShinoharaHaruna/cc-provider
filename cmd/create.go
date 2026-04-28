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

	// 2. Ask if user wants to use a template
	fmt.Println("\nWould you like to use a template? (y/n)")
	useTemplate := prompt(reader, "Use template", false)
	var envVars map[string]string

	if strings.ToLower(useTemplate) == "y" || strings.ToLower(useTemplate) == "yes" {
		// List available templates
		templates, err := listTemplates()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing templates: %v\n", err)
			os.Exit(1)
		}

		if len(templates) == 0 {
			fmt.Println("No templates available. Creating environment manually.")
			envVars = make(map[string]string)
		} else {
			fmt.Println("\nAvailable templates:")
			for i, tmpl := range templates {
				fmt.Printf("  %d. %s - %s\n", i+1, tmpl.Name, tmpl.Description)
			}

			selection := prompt(reader, "\nSelect template number (or press Enter to skip)", false)
			if selection != "" {
				var idx int
				if _, err := fmt.Sscanf(selection, "%d", &idx); err == nil && idx > 0 && idx <= len(templates) {
					selectedTmpl := templates[idx-1]
					envVars = make(map[string]string)
					for k, v := range selectedTmpl.EnvVars {
						envVars[k] = v
					}
					fmt.Printf("\nUsing template '%s'.\n", selectedTmpl.Name)
				} else {
					fmt.Println("Invalid selection. Creating environment manually.")
					envVars = make(map[string]string)
				}
			} else {
				envVars = make(map[string]string)
			}
		}
	} else {
		envVars = make(map[string]string)
	}

	// 3. Required inputs
	fmt.Println("\nEnter required variables: ")
	if baseURL := promptWithExisting(reader, "  ANTHROPIC_BASE_URL", envVars["ANTHROPIC_BASE_URL"], true); baseURL != "" {
		envVars["ANTHROPIC_BASE_URL"] = baseURL
	}
	if authToken := promptWithExisting(reader, "  ANTHROPIC_AUTH_TOKEN", envVars["ANTHROPIC_AUTH_TOKEN"], true); authToken != "" {
		envVars["ANTHROPIC_AUTH_TOKEN"] = authToken
	}

	// 4. Recommended inputs
	fmt.Println("\nEnter recommended variables (press Enter to skip): ")
	if model := promptWithExisting(reader, "  ANTHROPIC_MODEL", envVars["ANTHROPIC_MODEL"], false); model != "" {
		envVars["ANTHROPIC_MODEL"] = model
	}
	if haikuModel := promptWithExisting(reader, "  ANTHROPIC_DEFAULT_HAIKU_MODEL", envVars["ANTHROPIC_DEFAULT_HAIKU_MODEL"], false); haikuModel != "" {
		envVars["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = haikuModel
	}
	if sonnetModel := promptWithExisting(reader, "  ANTHROPIC_DEFAULT_SONNET_MODEL", envVars["ANTHROPIC_DEFAULT_SONNET_MODEL"], false); sonnetModel != "" {
		envVars["ANTHROPIC_DEFAULT_SONNET_MODEL"] = sonnetModel
	}
	if opusModel := promptWithExisting(reader, "  ANTHROPIC_DEFAULT_OPUS_MODEL", envVars["ANTHROPIC_DEFAULT_OPUS_MODEL"], false); opusModel != "" {
		envVars["ANTHROPIC_DEFAULT_OPUS_MODEL"] = opusModel
	}
	if subagentModel := promptWithExisting(reader, "  CLAUDE_CODE_SUBAGENT_MODEL", envVars["CLAUDE_CODE_SUBAGENT_MODEL"], false); subagentModel != "" {
		envVars["CLAUDE_CODE_SUBAGENT_MODEL"] = subagentModel
	}
	if effortLevel := promptWithExisting(reader, "  CLAUDE_CODE_EFFORT_LEVEL", envVars["CLAUDE_CODE_EFFORT_LEVEL"], false); effortLevel != "" {
		envVars["CLAUDE_CODE_EFFORT_LEVEL"] = effortLevel
	}

	// 5. Optional inputs with defaults
	fmt.Println("\nEnter optional variables (press Enter to use default): ")
	timeout := promptWithExisting(reader, "  API_TIMEOUT_MS", envVars["API_TIMEOUT_MS"], false)
	if timeout == "" {
		timeout = "600000"
	}
	envVars["API_TIMEOUT_MS"] = timeout

	disableTraffic := promptWithExisting(reader, "  CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", envVars["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"], false)
	if disableTraffic == "" {
		disableTraffic = "1"
	}
	envVars["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = disableTraffic

	// 6. Write to file
	var finalEnvVars []string
	for key, value := range envVars {
		finalEnvVars = append(finalEnvVars, fmt.Sprintf("%s=\"%s\"", key, value))
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
