package cmd

import (
	"fmt"

	"github.com/dorochadev/aka/launcher"
	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <shortname>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a launcher",
	Long:    `Remove an existing launcher by its shortname.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("force", "f", false, "Remove without confirmation")
}

func runRemove(cmd *cobra.Command, args []string) error {
	shortname := args[0]

	if !launcher.Exists(shortname) {
		ui.PrintError(fmt.Sprintf("Launcher '%s' does not exist", shortname))
		return fmt.Errorf("launcher not found")
	}

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		confirm := ui.Confirm(fmt.Sprintf("Remove launcher '%s'?", shortname))
		if !confirm {
			ui.PrintInfo("Cancelled.")
			return nil
		}
	}

	if err := launcher.Remove(shortname); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to remove launcher: %v", err))
		return err
	}

	fmt.Println()
	ui.SuccessBox(fmt.Sprintf("Removed launcher '%s'", shortname))
	fmt.Println()

	return nil
}
