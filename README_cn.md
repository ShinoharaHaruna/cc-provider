# Claude Code Provider (`cc-provider`)

[![CI](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/ci.yml/badge.svg)](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/ci.yml)
[![Release](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/release.yml/badge.svg)](https://github.com/ShinoharaHaruna/cc-provider/actions/workflows/release.yml)
[![License](https://img.shields.io/github/license/ShinoharaHaruna/cc-provider)](LICENSE)

[English](README.md)

一个命令行工具，用于管理 Claude Code 的不同环境变量集，让您轻松切换各种 API 提供商，如 DeepSeek、Anthropic 等。

## 核心功能

- **创建和管理环境**：交互式地创建、列出和移除提供商环境。
- **激活环境**：为您的 shell 会话持久设置环境变量。
- **Shell 集成**：自动挂接到 `.zshrc` 或 `.bashrc`，实现无缝激活。
- **导出配置**：将环境设置导出为 `.env` 文件格式。

## 安装

### 下载预构建二进制文件

从 [Releases](https://github.com/ShinoharaHaruna/cc-provider/releases) 页面下载适用于您平台的最新版本。

支持的平台：

- Linux
  - linux-amd64
  - linux-arm64
- macOS
  - darwin-amd64
  - darwin-arm64

```bash
# 下载并解压（根据需要替换平台）
wget https://github.com/ShinoharaHaruna/cc-provider/releases/latest/download/cc-provider-linux-amd64.tar.gz
tar -xzf cc-provider-linux-amd64.tar.gz

# 移动到 PATH 中的目录
sudo mv cc-provider /usr/local/bin/
```

### 从源代码构建

1. **先决条件**：您需要安装 Go (版本 1.18 或更高)。

2. **从源代码安装**：克隆仓库并使用版本信息构建。

    ```bash
    git clone https://github.com/ShinoharaHaruna/cc-provider.git
    cd cc-provider
    make
    make install
    ```

    这会编译二进制文件并将其放置在您的 `$GOPATH/bin` 目录中。为了全局运行 `cc-provider` 命令，请确保此目录在您系统的 `PATH` 中。

    您可以通过将以下行添加到您的 shell 配置文件（例如 `~/.zshrc` 或 `~/.bashrc`）来做到这一点：

    ```bash
    export PATH=$PATH:$(go env GOPATH)/bin
    ```

3. **本地构建**：要不安装而构建二进制文件：

    ```bash
    make build
    ```

    二进制文件将放置在 `bin/cc-provider` 中。

## 工作原理

该工具在 `~/.cc-provider/` 目录中管理环境配置文件。当您激活一个环境时，它会生成一个脚本 `~/.cc-provider/active_env.sh`。

### 自动设置

首次运行任何 `cc-provider` 命令时，该工具将自动执行一次性设置：

1. 它会将一个 `source` 命令附加到您的 shell 配置文件（例如 `~/.zshrc` 或 `~/.bashrc`）。
2. 它会为您的 shell 生成并安装一个 Tab 补全脚本。
3. 它会创建一个 shell 函数包装器，以实现即时激活。

这确保了在每个新的 shell 会话中都可以使用环境变量、Tab 补全和即时激活。您只需在初始设置后重启 shell 一次。

设置完成后，`cc-provider activate` 的工作方式就像 `conda activate` 一样——无需 `eval` 或重启 shell！

## 命令

### `cc-provider list`

列出所有可用的提供商环境。当前激活的环境用星号 (`*`) 标记。

```bash
cc-provider list
```

### `cc-provider create`

交互式地创建一个新的提供商环境。系统将提示您输入环境名称和所需/可选变量。

```bash
cc-provider create
```

### `cc-provider activate <env-name>`

立即在当前 shell 中激活指定环境（无需重启）。

```bash
cc-provider activate deepseek
```

### `cc-provider remove <env-name>`

移除指定环境。如果环境当前处于激活状态，它将被停用。

```bash
cc-provider remove deepseek
```

### `cc-provider export`

将当前激活环境的配置以 `.env` 格式导出到标准输出。

```bash
cc-provider export
```

### `cc-provider export --name <env-name>`

导出特定环境的配置。

```bash
cc-provider export --name deepseek
```

### `cc-provider modify [env-name]`

交互式地修改现有提供商环境。如果未提供环境名称，系统将提示您从可用环境中选择。

```bash
cc-provider modify deepseek

# 或交互式选择：
cc-provider modify
```

### `cc-provider version`

显示版本信息，包括语义版本、构建时间和 git 提交哈希。

```bash
cc-provider version
```
