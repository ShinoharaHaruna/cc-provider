# Claude Code Provider (`cc-provider`)

[![CI](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/ci.yml/badge.svg)](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/ci.yml)
[![Release](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/release.yml/badge.svg)](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/release.yml)
[![License](https://img.shields.io/github/license/ShinoharaHaruna/cc-provider)](LICENSE)

A command-line tool to manage different sets of environment variables for Claude Code, allowing you to easily switch between various API providers like DeepSeek, Anthropic, etc.

## Core Features

- **Create & Manage Environments**: Interactively create, list, and remove provider environments.
- **Activate Environments**: Persistently set environment variables for your shell sessions.
- **Shell Integration**: Automatically hooks into `.zshrc` or `.bashrc` for seamless activation.
- **Export Configurations**: Export environment settings to a `.env` file format.

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/ShinoharaHaruna/cc-provider/releases) page.

**Linux/macOS:**

```bash
# Download and extract (replace VERSION and PLATFORM as needed)
wget https://github.com/ShinoharaHaruna/cc-provider/releases/download/v0.1.0/cc-provider-linux-amd64.tar.gz
tar -xzf cc-provider-linux-amd64.tar.gz

# Move to a directory in your PATH
sudo mv cc-provider /usr/local/bin/
```

### Build from Source

1. **Prerequisites**: You need to have Go installed (version 1.18 or newer).

2. **Install from source**: Clone the repository and build with version information.

    ```bash
    git clone https://github.com/ShinoharaHaruna/cc-provider.git
    cd cc-provider
    make
    make install
    ```

    This will compile the binary and place it in your `$GOPATH/bin` directory. To run the `cc-provider` command globally, ensure that this directory is in your system's `PATH`.

    You can do this by adding the following line to your shell's configuration file (e.g., `~/.zshrc` or `~/.bashrc`):

    ```bash
    export PATH=$PATH:$(go env GOPATH)/bin
    ```

3. **Build locally**: To build the binary without installing:

    ```bash
    make build
    ```

    The binary will be placed in `bin/cc-provider`.

## How It Works

The tool manages environment configuration files in the `~/.cc-provider/` directory. When you activate an environment, it generates a script `~/.cc-provider/active_env.sh`.

### Automatic Setup

The first time you run any `cc-provider` command, the tool will automatically perform a one-time setup:

1. It appends a `source` command to your shell's configuration file (e.g., `~/.zshrc` or `~/.bashrc`).
2. It generates and installs a tab completion script for your shell.
3. It creates a shell function wrapper that enables immediate activation.

This ensures that environment variables, tab completion, and instant activation are available in every new shell session. You just need to restart your shell once after the initial setup.

After setup, `cc-provider activate` works just like `conda activate` - no need for `eval` or shell restart!

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

Activates the specified environment immediately in the current shell (no restart needed).

```bash
cc-provider activate deepseek
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

### `cc-provider modify [env-name]`

Interactively modifies an existing provider environment. If no environment name is provided, you will be prompted to select from available environments.

```bash
cc-provider modify deepseek

# Or select interactively:
cc-provider modify
```

### `cc-provider version`

Displays version information including the semantic version, build time, and git commit hash.

```bash
cc-provider version
```
