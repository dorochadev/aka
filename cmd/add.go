package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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
	Args: cobra.MinimumNArgs(2),
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
	// Collect all remaining args as potential targets
	rawTargets := args[1:]

	var targets []string
	if len(rawTargets) == 1 && strings.Contains(rawTargets[0], ",") {
		// Handle comma-separated list: aka add dev "VS Code,Safari"
		split := strings.Split(rawTargets[0], ",")
		for _, t := range split {
			targets = append(targets, strings.TrimSpace(t))
		}
	} else {
		// Handle space-separated args: aka add dev "VS Code" Safari
		targets = rawTargets
	}

	// Determine if it's a stack or single launcher
	isStack := len(targets) > 1
	var target string
	if !isStack {
		target = targets[0]
	}

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

	var launcherType launcher.LauncherType
	if isStack {
		launcherType = launcher.TypeStack
	} else {
		launcherType = launcher.DetectLauncherType(target)
	}

	metadata := &launcher.LauncherMetadata{
		Type:    launcherType,
		Target:  target,
		Targets: targets,
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
	case launcher.TypeStack:
		ui.SuccessBox(fmt.Sprintf("Created stack launcher '%s' with %d items", shortname, len(targets)))
		for _, t := range targets {
			ui.PrintResult("-", t)
		}
		ui.PrintExample("Launch all:", shortname)
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
		if ui.Confirm("Add launcher directory to your PATH?") {
			printReloadInstructions()
		} else {
			ui.PrintInfo("Run 'source ~/.zshrc' or restart your terminal to use the launcher.")
			fmt.Println()
		}
	}

	return nil
}

func printReloadInstructions() {
	shell := os.Getenv("SHELL")

	var configFile string
	switch {
	case strings.Contains(shell, "zsh"):
		configFile = "~/.zshrc"
	case strings.Contains(shell, "bash"):
		configFile = "~/.bashrc"
	default:
		configFile = "your shell config"
	}

	fmt.Println()
	ui.PrintInfo("To activate the launcher, run:")
	fmt.Println()
	ui.PrintCommand(fmt.Sprintf("source %s", configFile))
	fmt.Println()
}

func isValidShortname(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}
