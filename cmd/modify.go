package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify [env-name]",
	Short: "Interactively modifies an existing provider environment.",
	Long:  `Interactively prompts for the necessary details to modify an existing provider environment file in the ~/.cc-provider directory.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runModifyCmd,
}

func runModifyCmd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var envName string
	// 如果没有提供环境名称,列出可用环境供用户选择 / If no environment name is provided, list available environments for user selection
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

	// 验证环境是否存在 / Validate environment exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Environment '%s' not found.\n", envName)
		os.Exit(1)
	}

	// 读取现有环境变量 / Read existing environment variables
	existingVars, err := readEnvFile(envFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading environment file '%s': %v\n", envFilePath, err)
		os.Exit(1)
	}

	fmt.Printf("\nModifying environment '%s'...\n", envName)
	fmt.Println("Press Enter to keep current value, or enter new value to update.")

	var envVars []string

	// 必填项 / Required inputs
	fmt.Println("\nRequired variables: ")
	baseURL := promptWithExisting(reader, "  ANTHROPIC_BASE_URL", existingVars["ANTHROPIC_BASE_URL"], true)
	envVars = append(envVars, "ANTHROPIC_BASE_URL="+baseURL)

	authToken := promptWithExisting(reader, "  ANTHROPIC_AUTH_TOKEN", existingVars["ANTHROPIC_AUTH_TOKEN"], true)
	envVars = append(envVars, "ANTHROPIC_AUTH_TOKEN="+authToken)

	// 推荐项 / Recommended inputs
	fmt.Println("\nRecommended variables (press Enter to keep current or clear): ")
	model := promptWithExisting(reader, "  ANTHROPIC_MODEL", existingVars["ANTHROPIC_MODEL"], false)
	if model != "" {
		envVars = append(envVars, "ANTHROPIC_MODEL="+model)
	}

	smallModel := promptWithExisting(reader, "  ANTHROPIC_SMALL_FAST_MODEL", existingVars["ANTHROPIC_SMALL_FAST_MODEL"], false)
	if smallModel != "" {
		envVars = append(envVars, "ANTHROPIC_SMALL_FAST_MODEL="+smallModel)
	}

	// 可选项(带默认值) / Optional inputs with defaults
	fmt.Println("\nOptional variables: ")
	timeout := promptWithExisting(reader, "  API_TIMEOUT_MS", existingVars["API_TIMEOUT_MS"], false)
	if timeout == "" {
		timeout = "600000"
	}
	envVars = append(envVars, "API_TIMEOUT_MS="+timeout)

	disableTraffic := promptWithExisting(reader, "  CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", existingVars["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"], false)
	if disableTraffic == "" {
		disableTraffic = "1"
	}
	envVars = append(envVars, "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC="+disableTraffic)

	// 写入文件 / Write to file
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

	fmt.Printf("\nSuccessfully modified environment '%s'.\n", envName)

	// 检查是否是当前激活的环境 / Check if this is the currently active environment
	if activeEnv := os.Getenv("CC_PROVIDER_ACTIVE_ENV"); activeEnv == envName {
		fmt.Println("Note: This is the currently active environment.")
		fmt.Printf("To apply changes, run: cc-provider activate %s\n", envName)
	}
}

// selectEnvironment lists available environments and prompts user to select one
// 列出可用环境并提示用户选择
func selectEnvironment(reader *bufio.Reader) string {
	entries, err := os.ReadDir(cfgDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config directory: %v\n", err)
		os.Exit(1)
	}

	var envs []string
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != "active_env.sh" {
			envs = append(envs, entry.Name())
		}
	}

	if len(envs) == 0 {
		fmt.Println("No environments found. Use 'cc-provider create' to create one.")
		return ""
	}

	fmt.Println("\nAvailable environments:")
	for i, env := range envs {
		fmt.Printf("  %d. %s\n", i+1, env)
	}

	fmt.Print("\nSelect environment number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selection int
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil || selection < 1 || selection > len(envs) {
		fmt.Fprintf(os.Stderr, "Invalid selection.\n")
		return ""
	}

	return envs[selection-1]
}

// readEnvFile reads an environment file and returns a map of key-value pairs
// 读取环境文件并返回键值对映射
func readEnvFile(filePath string) (map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], `"`)
			envVars[key] = value
		}
	}

	return envVars, scanner.Err()
}

// promptWithExisting prompts user with existing value as default
// 提示用户输入,显示现有值作为默认值
func promptWithExisting(reader *bufio.Reader, message, existingValue string, required bool) string {
	if existingValue != "" {
		fmt.Printf("%s (current: %s): ", message, existingValue)
	} else {
		fmt.Printf("%s: ", message)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 如果用户输入为空 / If user input is empty
	if input == "" {
		// 如果是必填项且没有现有值,继续提示 / If required and no existing value, continue prompting
		if required && existingValue == "" {
			fmt.Println("This field is required.")
			return promptWithExisting(reader, message, existingValue, required)
		}
		// 返回现有值 / Return existing value
		return existingValue
	}

	return input
}

func init() {
	rootCmd.AddCommand(modifyCmd)
}
