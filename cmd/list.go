package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available provider environments.",
	Long:  `Lists all provider environments configured in the ~/.cc-provider directory. The active environment is marked with an asterisk (*).`,
	Run: func(cmd *cobra.Command, args []string) {
		activeEnv := os.Getenv("CC_PROVIDER_ACTIVE_ENV")

		files, err := os.ReadDir(cfgDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config directory '%s'\n", cfgDir)
			os.Exit(1)
		}

		var envs []string
		for _, file := range files {
			fileName := file.Name()
			// 过滤掉系统文件 / Filter out system files
			if !file.IsDir() &&
				fileName != filepath.Base(activeEnvFile) &&
				!strings.HasPrefix(fileName, "completion.") &&
				fileName != "shell_function.sh" {
				envs = append(envs, fileName)
			}
		}

		if len(envs) == 0 {
			fmt.Println("No provider environments found. Use 'cc-provider create' to add one.")
			return
		}

		fmt.Println("Available provider environments: ")
		for _, env := range envs {
			if env == activeEnv {
				fmt.Printf("  * %s\n", env)
			} else {
				fmt.Printf("    %s\n", env)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
