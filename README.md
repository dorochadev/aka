# aka

A simple CLI tool to create launchers for your favorite applications, URLs, SSH connections, and commands.

## What it does

`aka` creates shortcut commands in `~/bin` that launch applications, open URLs, connect to servers, or run shell commands. Instead of typing long commands or searching for apps, just type a short alias.

## Features

- **Application launchers** - Open GUI apps with a short command
- **URL launchers** - Open websites in your default browser
- **SSH launchers** - Connect to servers with saved credentials
- **Command launchers** - Create shortcuts for shell commands
- **Stack launchers** - Open multiple apps/URLs with one command
- **Environment variables** - Set env vars for any launcher
- **Shell completions** - Tab-complete launcher names and commands
- **Simple management** - Add, remove, rename, and list launchers

## Installation

```bash
git clone https://github.com/dorochadev/aka.git
cd aka
go build -o aka
sudo mv aka /usr/local/bin/
```

## Quick Start

### Application Launchers

```bash
aka add safari Safari
aka add code "VS Code"
safari                    # Opens Safari
code myproject.code       # Opens VS Code with file
```

### URL Launchers

```bash
aka add gh https://github.com
aka add youtube https://youtube.com
gh                        # Opens GitHub in browser
```

### SSH Launchers

```bash
aka add server user@192.168.1.1
aka add prod user@prod.com --save-password  # Prompts for password securely
server                    # Connects via SSH
```

### Command Launchers

```bash
aka add ll "ls -lah"
aka add ports "lsof -i -P"
ll                        # Runs ls -lah
```

### Stack Launchers

Open multiple apps or URLs with a single command:

```bash
# Comma-separated list
aka add dev "VS Code,iTerm,Safari"
aka add work "Slack,https://mail.google.com,Notion"

# Space-separated arguments
aka add stack Spotify Discord Chrome

dev                       # Opens all 3 apps at once
```

### Environment Variables

```bash
aka add dev "VS Code" --env "DEBUG=1,ENV=development"
dev                       # Opens VS Code with env vars set
```

## Shell Completions

Install completions with one command:

```bash
aka completion install
# Auto-detects your shell (zsh/bash)
# Installs completions automatically
# Restart terminal or source your shell config
```

Now you can tab-complete:

```bash
aka <TAB>                 # Shows: add, remove, list, etc.
aka remove <TAB>          # Shows your launcher names
```

## Commands

```bash
aka add <name> <target>              # Create a launcher
aka remove <name>                    # Remove a launcher
aka list                             # List all launchers
aka rename <old> <new>               # Rename a launcher
aka open <name> [files...]           # Open launcher with files
aka completion install               # Install shell completions
```

### Flags

```bash
--save-password          # Prompt for SSH password (secure)
--env key=value          # Set environment variables
--port <number>          # SSH port (default: 22)
--key <path>             # SSH key file
-f, --force              # Overwrite without confirmation
```

## Examples

```bash
# Create launchers
aka add safari Safari
aka add gh https://github.com/dorochadev
aka add server root@192.168.1.1 --port 2222
aka add deploy "ssh user@prod && cd /app && git pull"

# Use them
safari
gh
server
deploy

# Manage launchers
aka list
aka rename server prod-server
aka remove old-launcher
```

## How it works

`aka` creates executable shell scripts in `~/bin` that:

- Open applications with `open -a` (macOS) or `xdg-open` (Linux)
- Open URLs in your default browser
- Connect via SSH with optional `sshpass` for passwords
- Execute shell commands

Launcher configuration is stored in `~/.config/aka/launchers.json`.

## Requirements

- Go 1.21+ (for building)
- macOS or Linux
- `sshpass` (optional, for SSH password storage)

## License

MIT

## Setup

On first run, aka will prompt you to add `~/bin` to your PATH if needed.
