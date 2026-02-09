package cmd

import (
	"fmt"

	"github.com/dorochadev/aka/launcher"
	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:     "rename <oldname> <newname>",
	Aliases: []string{"mv"},
	Short:   "Rename a launcher",
	Long:    `Rename an existing launcher without changing its target application.`,
	Args:    cobra.ExactArgs(2),
	RunE:    runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	oldName := args[0]
	newName := args[1]

	// Validate new name
	if !isValidShortname(newName) {
		ui.PrintError("Invalid new name. Use only alphanumeric characters, hyphens, and underscores.")
		return fmt.Errorf("invalid shortname")
	}

	// Check if old launcher exists
	if !launcher.Exists(oldName) {
		ui.PrintError(fmt.Sprintf("Launcher '%s' does not exist", oldName))
		return fmt.Errorf("launcher not found")
	}

	// Check if new name already exists
	if launcher.Exists(newName) {
		ui.PrintError(fmt.Sprintf("Launcher '%s' already exists", newName))
		return fmt.Errorf("launcher already exists")
	}

	// Rename the launcher
	if err := launcher.Rename(oldName, newName); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to rename launcher: %v", err))
		return err
	}

	fmt.Println()
	ui.SuccessBox(fmt.Sprintf("Renamed launcher '%s' to '%s'", oldName, newName))
	fmt.Println()

	return nil
}
