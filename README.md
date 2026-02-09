# aka

Create short commands to launch GUI applications.

## Install

```bash
go install github.com/dorochadev/aka@latest
```

Or build from source:

```bash
git clone https://github.com/dorochadev/aka.git
cd aka
go build -o aka
sudo mv aka /usr/local/bin/
```

## Quick Start

```bash
# Create a launcher
aka add safari Safari
aka add code "Visual Studio Code"

# Use it
safari
code myproject.txt

# Manage launchers
aka list
aka rename safari web
aka remove web
```

## Commands

- `aka add <name> <app>` - Create a new launcher
- `aka list` - Show all launchers
- `aka remove <name>` - Delete a launcher
- `aka rename <old> <new>` - Rename a launcher

## Setup

On first run, aka will prompt you to add `~/bin` to your PATH if needed.
