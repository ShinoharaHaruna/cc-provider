package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:               "inspect [env-name]",
	Short:             "Inspects a provider environment's configuration.",
	Long:              `Displays the configuration of a provider environment. If no name is given, prompts you to select one interactively.`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completeEnvironmentNames,
	Run:               runInspectCmd,
}

func runInspectCmd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var envName string
	if len(args) == 0 {
		envName = selectEnvironment(reader)
		if envName == "" {
			fmt.Fprintf(os.Stderr, "Error: No environment selected.\n")
			os.Exit(1)
		}
	} else {
		envName = args[0]
	}

	envFilePath := filepath.Join(cfgDir, envName)
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Environment '%s' not found.\n", envName)
		os.Exit(1)
	}

	envVars, err := readEnvFile(envFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading environment '%s': %v\n", envName, err)
		os.Exit(1)
	}

	activeEnv := os.Getenv("CC_PROVIDER_ACTIVE_ENV")
	activeMarker := ""
	if envName == activeEnv {
		activeMarker = " (active)"
	}

	fmt.Printf("Environment: %s%s\n", envName, activeMarker)
	fmt.Println("---")
	for _, key := range envVarKeys {
		if key == "CC_PROVIDER_ACTIVE_ENV" {
			continue
		}
		val, ok := envVars[key]
		if ok {
			// Mask token value for safety
			if key == "ANTHROPIC_AUTH_TOKEN" {
				val = maskSecret(val)
			}
			fmt.Printf("  %-42s %s\n", key+":", val)
		} else {
			fmt.Printf("  %-42s %s\n", key+":", "(not set)")
		}
	}
}

// maskSecret shows only the first 6 and last 4 characters of a secret.
func maskSecret(s string) string {
	const visible = 6
	if len(s) <= visible+4 {
		return "****"
	}
	return s[:visible] + "****" + s[len(s)-4:]
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
