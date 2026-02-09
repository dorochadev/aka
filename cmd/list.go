package cmd

import (
	"fmt"

	"github.com/dorochadev/aka/launcher"
	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all launchers",
	Long:    `Display all configured launchers and their target applications.`,
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	launchers, err := launcher.List()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to list launchers: %v", err))
		return err
	}

	if len(launchers) == 0 {
		fmt.Println()
		ui.PrintInfo("No launchers configured yet.")
		fmt.Println()
		ui.PrintExample("Create your first launcher:", "aka add ag \"Adobe Acrobat\"")
		fmt.Println()
		return nil
	}

	headers := []string{"Command", ui.IconArrow, "Application"}
	rows := make([][]string, len(launchers))
	for i, l := range launchers {
		rows[i] = []string{l.Name, "", l.Target}
	}

	fmt.Println()
	ui.Table(headers, rows)
	fmt.Println()
	ui.PrintInfo(fmt.Sprintf("Total: %d launcher(s)", len(launchers)))
	fmt.Println()

	return nil
}
