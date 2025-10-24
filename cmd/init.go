package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Init performs the initial setup for cc-provider.
// It ensures the configuration directory and necessary files exist,
// and sets up shell integration if needed.
func Init() {
	// This function is called from main.go before any command is executed.
	setupConfigPaths()
	if _, err := ensureShellConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during initial shell setup: %v\n", err)
		os.Exit(1)
	}
}

// setupConfigPaths initializes the configuration directory path and creates it if necessary.
func setupConfigPaths() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	cfgDir = filepath.Join(home, ".cc-provider")
	activeEnvFile = filepath.Join(cfgDir, "active_env.sh")

	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory '%s': %v\n", cfgDir, err)
			os.Exit(1)
		}
	}
}

// ensureShellConfig checks and modifies the user's shell configuration file.
func ensureShellConfig() (string, error) {
	shell := os.Getenv("SHELL")
	var shellType string
	var rcFileName string

	if strings.Contains(shell, "zsh") {
		shellType = "zsh"
		rcFileName = ".zshrc"
	} else if strings.Contains(shell, "bash") {
		shellType = "bash"
		rcFileName = ".bashrc"
	} else {
		// For unsupported shells, we don't attempt to configure automatically.
		return "", nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	rcPath := filepath.Join(home, rcFileName)

	// Check if the source line for active_env.sh already exists
	file, err := os.Open(rcPath)
	needsSetup := true
	if err == nil {
		scanner := bufio.NewScanner(file)
		sourceLine := fmt.Sprintf("source \"%s\"", activeEnvFile)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), sourceLine) {
				needsSetup = false
				break
			}
		}
		file.Close()
	}

	if !needsSetup {
		return rcPath, nil // Already configured
	}

	// Create an empty active_env.sh if it doesn't exist, to prevent source errors on shell startup.
	if _, err := os.Stat(activeEnvFile); os.IsNotExist(err) {
		if err := os.WriteFile(activeEnvFile, []byte("# cc-provider active environment script\n"), 0644); err != nil {
			return "", fmt.Errorf("creating empty active_env.sh: %w", err)
		}
	}

	// Append the source lines to the rc file
	f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("opening %s for appending: %w", rcPath, err)
	}
	defer f.Close()

	// Generate and write completion script
	completionFilePath := filepath.Join(cfgDir, "completion."+shellType)
	completionFile, err := os.Create(completionFilePath)
	if err != nil {
		return "", fmt.Errorf("creating completion file: %w", err)
	}
	defer completionFile.Close()

	switch shellType {
	case "bash":
		rootCmd.GenBashCompletion(completionFile)
	case "zsh":
		rootCmd.GenZshCompletion(completionFile)
	}

	// Create shell function file for cc-provider command wrapper
	// 创建 cc-provider 命令包装的 shell 函数文件
	shellFunctionPath := filepath.Join(cfgDir, "shell_function.sh")
	shellFunctionContent := `# cc-provider shell integration
# This wraps the cc-provider command to enable immediate activation

cc-provider() {
    local cmd="$1"
    shift
    
    if [ "$cmd" = "activate" ]; then
        # Check for --eval flag
        local use_eval=false
        local env_name=""
        
        for arg in "$@"; do
            case "$arg" in
                --eval|-e)
                    use_eval=true
                    ;;
                *)
                    env_name="$arg"
                    ;;
            esac
        done
        
        if [ "$use_eval" = "true" ]; then
            # Direct eval mode
            eval "$(command cc-provider activate --eval "$env_name")"
        else
            # Automatic eval mode for immediate activation
            eval "$(command cc-provider activate --eval "$env_name")"
        fi
    else
        # For all other commands, call the actual binary
        command cc-provider "$cmd" "$@"
    fi
}
`
	if err := os.WriteFile(shellFunctionPath, []byte(shellFunctionContent), 0644); err != nil {
		return "", fmt.Errorf("creating shell function file: %w", err)
	}

	// Add source commands for activation, completion, and shell function scripts
	// 添加激活、补全和 shell 函数脚本的 source 命令
	configBlock := fmt.Sprintf(`
# Added by cc-provider for environment activation and auto-completion
source "%s"
source "%s"
source "%s"
`, activeEnvFile, completionFilePath, shellFunctionPath)

	if _, err := f.WriteString(configBlock); err != nil {
		return "", fmt.Errorf("writing to %s: %w", rcPath, err)
	}

	fmt.Printf("First-time setup complete. Added configuration to '%s' for automatic environment loading and tab completion.\n", rcPath)
	fmt.Println("Please restart your shell or source your config file to apply changes.")
	return rcPath, nil
}
