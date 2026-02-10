package cmd

import (
	"fmt"
	"strings"

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

	metadata, _ := launcher.LoadMetadata()

	headers := []string{"Command", "Type", ui.IconArrow, "Target"}
	rows := make([][]string, len(launchers))
	for i, l := range launchers {
		launcherType := "app"
		displayTarget := l.Target

		if meta, ok := metadata[l.Name]; ok && meta != nil {
			launcherType = string(meta.Type)

			// For stacks, show the list of targets
			if meta.Type == launcher.TypeStack && len(meta.Targets) > 0 {
				targetShort := strings.Join(meta.Targets, ", ")
				if len(targetShort) > 50 {
					targetShort = targetShort[:47] + "..."
				}
				displayTarget = targetShort
			}
		}

		rows[i] = []string{l.Name, launcherType, "", displayTarget}
	}

	fmt.Println()
	ui.Table(headers, rows)
	fmt.Println()
	ui.PrintInfo(fmt.Sprintf("Total: %d launcher(s)", len(launchers)))
	fmt.Println()

	return nil
}
