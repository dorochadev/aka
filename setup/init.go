package setup

import (
	"fmt"

	"github.com/dorochadev/aka/launcher"
	"github.com/dorochadev/aka/ui"
)

var hasShownWelcome = false

// EnsureSetup performs first-run initialization
func EnsureSetup() error {
	if err := launcher.EnsureLauncherDir(); err != nil {
		return fmt.Errorf("failed to setup launcher directory: %w", err)
	}

	if !launcher.IsInPath() && !hasShownWelcome {
		showPathWarning()
	}

	return nil
}

// showPathWarning displays a warning if the launcher directory is not in PATH
func showPathWarning() {
	hasShownWelcome = true
	dir := launcher.GetLauncherDir()

	fmt.Println()
	ui.WarningBox(fmt.Sprintf("Launcher directory is not in your PATH: %s", dir))
	fmt.Println()
	ui.PrintInfo("To use launchers directly, add this directory to your PATH:")
	fmt.Println()
	ui.PrintCommand(fmt.Sprintf(`echo 'export PATH="%s:$PATH"' >> ~/.zshrc`, dir))
	ui.PrintCommand("source ~/.zshrc")
	fmt.Println()
}
