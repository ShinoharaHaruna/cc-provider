package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage provider configuration templates.",
	Long:  `List, add, or remove provider configuration templates.`,
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates.",
	Long:  `List all built-in and custom provider configuration templates.`,
	Run:   runTemplateListCmd,
}

var templateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a custom template.",
	Long:  `Interactively create a new custom provider configuration template.`,
	Run:   runTemplateAddCmd,
}

var templateRemoveCmd = &cobra.Command{
	Use:   "remove [template-name]",
	Short: "Remove a custom template.",
	Long:  `Remove a custom provider configuration template. Built-in templates cannot be removed.`,
	Args:  cobra.ExactArgs(1),
	Run:   runTemplateRemoveCmd,
}

var templateShowCmd = &cobra.Command{
	Use:   "show [template-name]",
	Short: "Show template details.",
	Long:  `Display the details of a specific template.`,
	Args:  cobra.ExactArgs(1),
	Run:   runTemplateShowCmd,
}

func runTemplateListCmd(cmd *cobra.Command, args []string) {
	templates, err := listTemplates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing templates: %v\n", err)
		os.Exit(1)
	}

	if len(templates) == 0 {
		fmt.Println("No templates available.")
		return
	}

	fmt.Println("Available templates:")
	for _, tmpl := range templates {
		isBuiltIn := ""
		if _, ok := builtInTemplates[tmpl.Name]; ok {
			isBuiltIn = " (built-in)"
		}
		fmt.Printf("  - %s%s: %s\n", tmpl.Name, isBuiltIn, tmpl.Description)
	}
}

func runTemplateAddCmd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	// Get template name
	name := prompt(reader, "Enter template name", true)
	if _, ok := builtInTemplates[name]; ok {
		fmt.Fprintf(os.Stderr, "Error: Template name '%s' conflicts with built-in template.\n", name)
		os.Exit(1)
	}

	// Check if custom template already exists
	tmpl, _ := getTemplate(name)
	if tmpl != nil {
		fmt.Fprintf(os.Stderr, "Error: Template '%s' already exists.\n", name)
		os.Exit(1)
	}

	// Get description
	description := prompt(reader, "Enter template description", true)

	// Get environment variables
	envVars := make(map[string]string)
	fmt.Println("\nEnter environment variables (press Enter to skip):")

	baseURL := prompt(reader, "  ANTHROPIC_BASE_URL", false)
	if baseURL != "" {
		envVars["ANTHROPIC_BASE_URL"] = baseURL
	}

	authToken := prompt(reader, "  ANTHROPIC_AUTH_TOKEN", false)
	if authToken != "" {
		envVars["ANTHROPIC_AUTH_TOKEN"] = authToken
	}

	model := prompt(reader, "  ANTHROPIC_MODEL", false)
	if model != "" {
		envVars["ANTHROPIC_MODEL"] = model
	}

	smallModel := prompt(reader, "  ANTHROPIC_SMALL_FAST_MODEL", false)
	if smallModel != "" {
		envVars["ANTHROPIC_SMALL_FAST_MODEL"] = smallModel
	}

	haikuModel := prompt(reader, "  ANTHROPIC_DEFAULT_HAIKU_MODEL", false)
	if haikuModel != "" {
		envVars["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = haikuModel
	}

	sonnetModel := prompt(reader, "  ANTHROPIC_DEFAULT_SONNET_MODEL", false)
	if sonnetModel != "" {
		envVars["ANTHROPIC_DEFAULT_SONNET_MODEL"] = sonnetModel
	}

	opusModel := prompt(reader, "  ANTHROPIC_DEFAULT_OPUS_MODEL", false)
	if opusModel != "" {
		envVars["ANTHROPIC_DEFAULT_OPUS_MODEL"] = opusModel
	}

	subagentModel := prompt(reader, "  CLAUDE_CODE_SUBAGENT_MODEL", false)
	if subagentModel != "" {
		envVars["CLAUDE_CODE_SUBAGENT_MODEL"] = subagentModel
	}

	effortLevel := prompt(reader, "  CLAUDE_CODE_EFFORT_LEVEL", false)
	if effortLevel != "" {
		envVars["CLAUDE_CODE_EFFORT_LEVEL"] = effortLevel
	}

	timeout := prompt(reader, "  API_TIMEOUT_MS", false)
	if timeout != "" {
		envVars["API_TIMEOUT_MS"] = timeout
	}

	disableTraffic := prompt(reader, "  CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", false)
	if disableTraffic != "" {
		envVars["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = disableTraffic
	}

	// Create template
	newTemplate := Template{
		Name:        name,
		Description: description,
		EnvVars:     envVars,
	}

	if err := saveCustomTemplate(newTemplate); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully created template '%s'.\n", name)
	fmt.Printf("You can now use this template when creating a new environment.\n")
}

func runTemplateRemoveCmd(cmd *cobra.Command, args []string) {
	name := args[0]

	if _, ok := builtInTemplates[name]; ok {
		fmt.Fprintf(os.Stderr, "Error: Cannot remove built-in template '%s'.\n", name)
		os.Exit(1)
	}

	if err := deleteCustomTemplate(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully removed template '%s'.\n", name)
}

func runTemplateShowCmd(cmd *cobra.Command, args []string) {
	name := args[0]

	tmpl, err := getTemplate(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	isBuiltIn := ""
	if _, ok := builtInTemplates[name]; ok {
		isBuiltIn = " (built-in)"
	}

	fmt.Printf("Template: %s%s\n", tmpl.Name, isBuiltIn)
	fmt.Printf("Description: %s\n", tmpl.Description)
	fmt.Println("\nEnvironment variables:")
	for key, value := range tmpl.EnvVars {
		if key == "ANTHROPIC_AUTH_TOKEN" {
			value = maskSecret(value)
		}
		fmt.Printf("  %s=\"%s\"\n", key, value)
	}
}

func completeTemplateNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	templates, err := listTemplates()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, tmpl := range templates {
		names = append(names, tmpl.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateAddCmd)
	templateCmd.AddCommand(templateRemoveCmd)
	templateCmd.AddCommand(templateShowCmd)

	templateRemoveCmd.ValidArgsFunction = completeTemplateNames
	templateShowCmd.ValidArgsFunction = completeTemplateNames
}
