package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Manage shell completions",
	Long:  `Install, uninstall, or generate shell completion scripts for aka.`,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install shell completions automatically",
	Long:  `Automatically detect your shell and install completions.`,
	RunE:  runInstall,
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall shell completions",
	Long:  `Remove installed completion scripts and clean up shell configuration.`,
	RunE:  runUninstall,
}

var generateCmd = &cobra.Command{
	Use:   "generate [bash|zsh|fish]",
	Short: "Generate completion script",
	Long:  `Generate a completion script for the specified shell.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGenerate,
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(installCmd)
	completionCmd.AddCommand(uninstallCmd)
	completionCmd.AddCommand(generateCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	shell := detectShell()

	if shell == "unknown" {
		ui.PrintError("Could not detect shell. Please specify manually with 'aka completion generate [shell]'")
		return fmt.Errorf("unknown shell")
	}

	fmt.Println()
	if !ui.Confirm(fmt.Sprintf("Install completions for %s?", shell)) {
		ui.PrintInfo("Cancelled.")
		return nil
	}

	var err error
	switch shell {
	case "zsh":
		err = installZsh()
	case "bash":
		err = installBash()
	case "fish":
		ui.PrintError("Fish shell support coming soon. Use 'aka completion generate fish' for manual installation.")
		return fmt.Errorf("fish not yet supported")
	}

	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to install completions: %v", err))
		return err
	}

	return nil
}

func runUninstall(cmd *cobra.Command, args []string) error {
	shell := detectShell()

	if shell == "unknown" {
		ui.PrintError("Could not detect shell.")
		return fmt.Errorf("unknown shell")
	}

	fmt.Println()
	if !ui.Confirm(fmt.Sprintf("Uninstall completions for %s?", shell)) {
		ui.PrintInfo("Cancelled.")
		return nil
	}

	var err error
	switch shell {
	case "zsh":
		err = uninstallZsh()
	case "bash":
		err = uninstallBash()
	}

	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to uninstall completions: %v", err))
		return err
	}

	return nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	shell := args[0]

	var script string
	switch shell {
	case "zsh":
		script = generateZshCompletion()
	case "bash":
		script = generateBashCompletion()
	case "fish":
		ui.PrintError("Fish completions not yet implemented")
		return fmt.Errorf("fish not supported")
	default:
		ui.PrintError(fmt.Sprintf("Unknown shell: %s", shell))
		return fmt.Errorf("unknown shell")
	}

	fmt.Println(script)
	return nil
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "bash") {
		return "bash"
	}
	if strings.Contains(shell, "fish") {
		return "fish"
	}
	return "unknown"
}

func installZsh() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	completionDir := filepath.Join(home, ".zsh", "completions")
	if err := os.MkdirAll(completionDir, 0755); err != nil {
		return err
	}

	completionFile := filepath.Join(completionDir, "_aka")
	script := generateZshCompletion()

	if err := os.WriteFile(completionFile, []byte(script), 0644); err != nil {
		return err
	}

	zshrc := filepath.Join(home, ".zshrc")
	content, _ := os.ReadFile(zshrc)

	fpath := "fpath=(~/.zsh/completions $fpath)"
	autoload := "autoload -Uz compinit && compinit"

	if !strings.Contains(string(content), "~/.zsh/completions") {
		addition := "\n# aka completions\n" + fpath + "\n" + autoload + "\n"
		if err := appendToFile(zshrc, addition); err != nil {
			return err
		}
	}

	fmt.Println()
	ui.SuccessBox("Zsh completions installed!")
	fmt.Println()
	ui.PrintInfo("Restart your terminal or run:")
	ui.PrintCommand("source ~/.zshrc")
	fmt.Println()

	return nil
}

func installBash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	completionDir := filepath.Join(home, ".bash_completion.d")
	if err := os.MkdirAll(completionDir, 0755); err != nil {
		return err
	}

	completionFile := filepath.Join(completionDir, "aka")
	script := generateBashCompletion()

	if err := os.WriteFile(completionFile, []byte(script), 0644); err != nil {
		return err
	}

	bashrc := filepath.Join(home, ".bashrc")
	content, _ := os.ReadFile(bashrc)

	source := "source ~/.bash_completion.d/aka"
	if !strings.Contains(string(content), source) {
		addition := "\n# aka completions\n" + source + "\n"
		if err := appendToFile(bashrc, addition); err != nil {
			return err
		}
	}

	fmt.Println()
	ui.SuccessBox("Bash completions installed!")
	fmt.Println()
	ui.PrintInfo("Restart your terminal or run:")
	ui.PrintCommand("source ~/.bashrc")
	fmt.Println()

	return nil
}

func uninstallZsh() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	completionFile := filepath.Join(home, ".zsh", "completions", "_aka")
	os.Remove(completionFile)

	fmt.Println()
	ui.SuccessBox("Zsh completions uninstalled!")
	fmt.Println()
	ui.PrintInfo("You may want to remove the completion setup from ~/.zshrc manually.")
	fmt.Println()

	return nil
}

func uninstallBash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	completionFile := filepath.Join(home, ".bash_completion.d", "aka")
	os.Remove(completionFile)

	fmt.Println()
	ui.SuccessBox("Bash completions uninstalled!")
	fmt.Println()
	ui.PrintInfo("You may want to remove the completion setup from ~/.bashrc manually.")
	fmt.Println()

	return nil
}

func generateZshCompletion() string {
	return `#compdef aka

_aka() {
    local -a commands launchers
    
    commands=(
        'add:Create a new launcher'
        'remove:Remove a launcher'
        'list:List all launchers'
        'rename:Rename a launcher'
        'open:Open an application'
        'completion:Manage shell completions'
    )
    
    # Get launcher names from ~/bin
    if [[ -d ~/bin ]]; then
        launchers=(${(f)"$(ls ~/bin 2>/dev/null | grep -v '^\.')"})
    fi
    
    _arguments \
        '1: :->command' \
        '*::arg:->args'
    
    case $state in
        command)
            _describe 'command' commands
            ;;
        args)
            case $words[1] in
                remove|rename|open)
                    _describe 'launcher' launchers
                    ;;
                completion)
                    local -a subcommands
                    subcommands=(
                        'install:Install completions'
                        'uninstall:Uninstall completions'
                        'generate:Generate completion script'
                    )
                    _describe 'subcommand' subcommands
                    ;;
            esac
            ;;
    esac
}

_aka
`
}

func generateBashCompletion() string {
	return `_aka_completion() {
    local cur prev commands launchers
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    commands="add remove list rename open completion"
    
    if [ -d ~/bin ]; then
        launchers=$(ls ~/bin 2>/dev/null | grep -v '^\.')
    fi
    
    case "${prev}" in
        remove|rename|open)
            COMPREPLY=($(compgen -W "${launchers}" -- ${cur}))
            return 0
            ;;
        completion)
            COMPREPLY=($(compgen -W "install uninstall generate" -- ${cur}))
            return 0
            ;;
        aka)
            COMPREPLY=($(compgen -W "${commands}" -- ${cur}))
            return 0
            ;;
    esac
}

complete -F _aka_completion aka
`
}

func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}
