package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var exportName string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports an environment's configuration in .env format.",
	Long: `Exports the configuration of a specified environment to standard output in KEY=VALUE format.
If no environment is specified with the --name flag, it defaults to the currently active environment.`,
	Run: runExportCmd,
}

func runExportCmd(cmd *cobra.Command, args []string) {
	envName := exportName

	// If --name is not provided, use the active environment
	if envName == "" {
		envName = os.Getenv("CC_PROVIDER_ACTIVE_ENV")
		if envName == "" {
			fmt.Fprintln(os.Stderr, "Error: No environment specified and no active environment found.")
			fmt.Fprintln(os.Stderr, "Use '--name <env-name>' or activate an environment first.")
			os.Exit(1)
		}
	}

	envFilePath := filepath.Join(cfgDir, envName)

	// Validate environment exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Environment '%s' not found.\n", envName)
		os.Exit(1)
	}

	// Read and print the environment file content
	content, err := os.ReadFile(envFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading environment file '%s': %v\n", envFilePath, err)
		os.Exit(1)
	}

	// The content is already in KEY="VALUE" format, so just print it.
	fmt.Print(string(content))
	// Ensure there's a newline at the end of the output.
	if !strings.HasSuffix(string(content), "\n") {
		fmt.Println()
	}
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportName, "name", "", "Name of the environment to export")
}
