package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	activateEval bool // 是否输出 eval 格式 / Whether to output eval format
)

var activateCmd = &cobra.Command{
	Use:   "activate [env-name]",
	Short: "Activates a provider environment.",
	Long: `Activates a specified provider environment by updating the active environment script.

For immediate activation in current shell, use:
  eval "$(cc-provider activate --eval <env-name>)"

Otherwise, restart your shell or source your shell config file to apply changes.`,
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

	// 3. If --eval flag is set, output shell commands for immediate activation
	// 如果设置了 --eval 标志,输出 shell 命令以立即激活
	if activateEval {
		if err := outputEvalCommands(envName, envFilePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating eval commands: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 4. Normal mode: update config file and prompt user
	// 普通模式:更新配置文件并提示用户
	fmt.Printf("Successfully updated environment '%s' configuration.\n", envName)
	fmt.Println("\nThe environment will be active in new shell sessions.")
	fmt.Println("To activate immediately in current shell, run:")
	fmt.Printf("  cc-provider activate %s\n", envName)
	fmt.Printf("\n(If the shell function is not loaded, use: eval \"$(command cc-provider activate --eval %s)\")\n", envName)
}

// outputEvalCommands outputs shell commands for immediate activation via eval
// 输出用于通过 eval 立即激活的 shell 命令
func outputEvalCommands(envName, envFilePath string) error {
	var sb strings.Builder

	// Unset previous variables
	// 取消设置之前的变量
	for _, key := range envVarKeys {
		sb.WriteString(fmt.Sprintf("unset %s; ", key))
	}

	// Read variables from the environment file
	// 从环境文件读取变量
	content, err := os.ReadFile(envFilePath)
	if err != nil {
		return fmt.Errorf("reading env file %s: %w", envFilePath, err)
	}

	// Export new environment variables
	// 导出新的环境变量
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			sb.WriteString(fmt.Sprintf("export %s=%s; ", key, value))
		}
	}

	// Set the active environment identifier
	// 设置激活的环境标识符
	sb.WriteString(fmt.Sprintf(`export CC_PROVIDER_ACTIVE_ENV="%s"; `, envName))

	// Output success message to stderr so it doesn't interfere with eval
	// 将成功消息输出到 stderr,以免干扰 eval
	sb.WriteString(fmt.Sprintf(`echo "Environment '%s' activated." >&2`, envName))

	// Output the commands to stdout for eval
	// 将命令输出到 stdout 供 eval 使用
	fmt.Print(sb.String())
	return nil
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
	activateCmd.Flags().BoolVarP(&activateEval, "eval", "e", false, "Output shell commands for eval (use with: eval \"$(cc-provider activate --eval <env>)\")")
}
