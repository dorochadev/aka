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
	Use:   "add <shortname> <target>",
	Short: "Create a new launcher",
	Long: `Create a new launcher for an application, URL, SSH connection, or command.

The shortname should be alphanumeric and will become the command you type.
The target can be:
  - Application name (e.g., "Safari", "VS Code")
  - URL (e.g., https://youtube.com)
  - SSH connection (e.g., user@host)
  - Shell command (e.g., "ls -la")`,
	Args: cobra.ExactArgs(2),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("force", "f", false, "Overwrite existing launcher without confirmation")
	addCmd.Flags().Bool("save-password", false, "Prompt to save SSH password securely")
	addCmd.Flags().StringToString("env", nil, "Environment variables (key=value)")
	addCmd.Flags().IntP("port", "", 22, "SSH port")
	addCmd.Flags().StringP("key", "k", "", "SSH key file path")
}

func runAdd(cmd *cobra.Command, args []string) error {
	shortname := args[0]
	target := args[1]

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

	launcherType := launcher.DetectLauncherType(target)

	metadata := &launcher.LauncherMetadata{
		Type:   launcherType,
		Target: target,
	}

	envVars, _ := cmd.Flags().GetStringToString("env")
	if len(envVars) > 0 {
		metadata.Env = envVars
	}

	if launcherType == launcher.TypeSSH {
		savePassword, _ := cmd.Flags().GetBool("save-password")
		port, _ := cmd.Flags().GetInt("port")
		keyFile, _ := cmd.Flags().GetString("key")

		metadata.SSHConfig = &launcher.SSHConfig{
			Port:    port,
			KeyFile: keyFile,
		}

		if savePassword {
			password, err := ui.PromptPassword(fmt.Sprintf("🔒 Enter SSH password for %s (will be stored securely): ", target))
			if err != nil {
				ui.PrintError(fmt.Sprintf("Failed to read password: %v", err))
				return err
			}
			metadata.SSHConfig.Password = password
		}
	}

	if err := launcher.Create(shortname, metadata); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to create launcher: %v", err))
		return err
	}

	fmt.Println()

	switch launcherType {
	case launcher.TypeURL:
		ui.SuccessBox(fmt.Sprintf("Created URL launcher '%s' for %s", shortname, target))
		ui.PrintExample("Open the URL:", shortname)
	case launcher.TypeSSH:
		ui.SuccessBox(fmt.Sprintf("Created SSH launcher '%s' for %s", shortname, target))
		ui.PrintExample("Connect via SSH:", shortname)
	case launcher.TypeCommand:
		ui.SuccessBox(fmt.Sprintf("Created command launcher '%s'", shortname))
		ui.PrintExample("Run the command:", shortname)
	default:
		ui.SuccessBox(fmt.Sprintf("Created launcher '%s' for %s", shortname, target))
		ui.PrintExample("Open the application:", shortname)
		ui.PrintExample("Open with a file:", fmt.Sprintf("%s document.pdf", shortname))
	}

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

func isValidShortname(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}
