package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dorochadev/aka/launcher"
	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <shortname> [files...]",
	Short: "Open an application via its launcher",
	Long: `Convenience command to invoke a launcher through aka.

This is optional - you can run launchers directly once they're in your PATH.`,
	Args: cobra.MinimumNArgs(1),
	RunE: runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
	shortname := args[0]
	files := args[1:]

	if !launcher.Exists(shortname) {
		ui.PrintError(fmt.Sprintf("Launcher '%s' does not exist", shortname))
		return fmt.Errorf("launcher not found")
	}

	path := filepath.Join(launcher.GetLauncherDir(), shortname)

	execCmd := exec.Command(path, files...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to execute launcher: %v", err))
		return err
	}

	return nil
}
