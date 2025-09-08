# Claude Code Provider (`cc-provider`)

A command-line tool to manage different sets of environment variables for Claude Code, allowing you to easily switch between various API providers like DeepSeek, Anthropic, etc.

## Core Features

- **Create & Manage Environments**: Interactively create, list, and remove provider environments.
- **Activate Environments**: Persistently set environment variables for your shell sessions.
- **Shell Integration**: Automatically hooks into `.zshrc` or `.bashrc` for seamless activation.
- **Export Configurations**: Export environment settings to a `.env` file format.

## Installation

1. **Prerequisites**: You need to have Go installed (version 1.18 or newer).

2. **Install from source**: Clone the repository and use `go install`.

    ```bash
    git clone https://github.com/ShinoharaHaruna/cc-provider.git
    cd cc-provider
    go install .
    ```

    This will compile the binary and place it in your `$GOPATH/bin` directory. To run the `cc-provider` command globally, ensure that this directory is in your system's `PATH`.

    You can do this by adding the following line to your shell's configuration file (e.g., `~/.zshrc` or `~/.bashrc`):

    ```bash
    export PATH=$PATH:$(go env GOPATH)/bin
    ```

## How It Works

The tool manages environment configuration files in the `~/.cc-provider/` directory. When you activate an environment, it generates a script `~/.cc-provider/active_env.sh`.

### Automatic Setup

The first time you run any `cc-provider` command, the tool will automatically perform a one-time setup:

1. It appends a `source` command to your shell's configuration file (e.g., `~/.zshrc` or `~/.bashrc`).
2. It generates and installs a tab completion script for your shell.

This ensures that both environment variables and tab completion are automatically loaded in every new shell session. You just need to restart your shell once after the initial setup.

## Commands

### `cc-provider list`

Lists all available provider environments. The currently active environment is marked with an asterisk (`*`).

```bash
cc-provider list
```

### `cc-provider create`

Interactively creates a new provider environment. You will be prompted to enter the environment name and the required/optional variables.

```bash
cc-provider create
```

### `cc-provider activate <env-name>`

Activates the specified environment. After activation, you need to reload your shell's configuration or open a new terminal for the changes to take effect.

```bash
cc-provider activate deepseek

# Then, apply the changes:
source ~/.zshrc # Or ~/.bashrc
```

### `cc-provider remove <env-name>`

Removes the specified environment. If the environment is currently active, it will be deactivated.

```bash
cc-provider remove deepseek
```

### `cc-provider export`

Exports the currently active environment's configuration to standard output in `.env` format.

```bash
cc-provider export
```

### `cc-provider export --name <env-name>`

Exports a specific environment's configuration.

```bash
cc-provider export --name deepseek
```
