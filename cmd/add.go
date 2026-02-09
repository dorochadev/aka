package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/dorochadev/aka/launcher"
	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <shortname> <application>",
	Short: "Create a new launcher",
	Long: `Create a new launcher that opens the specified application.

The shortname should be alphanumeric and will become the command you type.
The application name should match the name of the GUI application.`,
	Args: cobra.ExactArgs(2),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("force", "f", false, "Overwrite existing launcher without confirmation")
}

func runAdd(cmd *cobra.Command, args []string) error {
	shortname := args[0]
	appName := args[1]

	if !isValidShortname(shortname) {
		ui.PrintError("Invalid shortname. Use only alphanumeric characters, hyphens, and underscores.")
		return fmt.Errorf("invalid shortname")
	}

	if launcher.Exists(shortname) {
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			overwrite := ui.Confirm(fmt.Sprintf("Launcher '%s' already exists. Overwrite?", shortname))
			if !overwrite {
				ui.PrintInfo("Cancelled.")
				return nil
			}
		}
	}

	if err := launcher.Create(shortname, appName); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to create launcher: %v", err))
		return err
	}

	fmt.Println()
	ui.SuccessBox(fmt.Sprintf("Created launcher '%s' for %s", shortname, appName))

	ui.PrintExample("Open the application:", shortname)
	ui.PrintExample("Open with a file:", fmt.Sprintf("%s document.pdf", shortname))
	fmt.Println()

	if !launcher.IsInPath() {
		if ui.Confirm("Reload your shell to use the launcher now?") {
			reloadShell()
		} else {
			ui.PrintInfo("Run 'source ~/.zshrc' or restart your terminal to use the launcher.")
			fmt.Println()
		}
	}

	return nil
}

// reloadShell provides instructions to reload the shell configuration
func reloadShell() {
	shell := os.Getenv("SHELL")

	var configFile string
	switch shell {
	case "/bin/zsh", "/usr/bin/zsh":
		configFile = "~/.zshrc"
	case "/bin/bash", "/usr/bin/bash":
		configFile = "~/.bashrc"
	default:
		configFile = "your shell config"
	}

	fmt.Println()
	ui.PrintInfo("To activate the launcher in your current shell, run:")
	fmt.Println()
	ui.PrintCommand(fmt.Sprintf("source %s", configFile))
	fmt.Println()
	ui.PrintInfo("Or simply restart your terminal.")
	fmt.Println()
}

// isValidShortname checks if a shortname contains only valid characters
func isValidShortname(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}
