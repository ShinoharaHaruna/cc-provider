package cmd

import (
	"os"
	"strings"

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

// getEnvironmentNames returns a list of available environment names for completion
// 返回可用环境名称列表用于补全
func getEnvironmentNames() []string {
	files, err := os.ReadDir(cfgDir)
	if err != nil {
		return []string{}
	}

	var envs []string
	for _, file := range files {
		fileName := file.Name()
		// 过滤掉系统文件 / Filter out system files
		if !file.IsDir() &&
			fileName != "active_env.sh" &&
			!strings.HasPrefix(fileName, "completion.") &&
			fileName != "shell_function.sh" {
			envs = append(envs, fileName)
		}
	}
	return envs
}
