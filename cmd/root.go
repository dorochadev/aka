package cmd

import (
	"fmt"
	"os"

	"github.com/dorochadev/aka/setup"
	"github.com/dorochadev/aka/ui"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "aka",
	Short: "Create short commands for launching GUI applications",
	Long: `aka is a CLI tool that generates executable launcher files for GUI applications.

Instead of storing shell aliases, aka creates small executable files that
open your applications. The filesystem itself acts as the database.`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return setup.EnsureSetup()
	},
	SilenceErrors: true, // We'll handle errors ourselves
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
}

func init() {
	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Set custom help command for root
	rootCmd.SetHelpFunc(customHelp)

	// Set custom usage template for all commands
	cobra.AddTemplateFunc("styleUsage", styleUsage)
	cobra.AddTemplateFunc("styleFlags", styleFlags)
	rootCmd.SetUsageTemplate(getUsageTemplate())

	// Set custom help function for all subcommands
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show help for command")
}

// styleUsage returns styled usage text
func styleUsage(s string) string {
	return ui.CurrentTheme.Body.Sprint(s)
}

// styleFlags returns styled flag text
func styleFlags(s string) string {
	return ui.CurrentTheme.Body.Sprint(s)
}

// customHelp displays a themed help menu for root command
func customHelp(cmd *cobra.Command, args []string) {
	// If this is a subcommand, show subcommand help
	if cmd.Parent() != nil {
		subcommandHelp(cmd)
		return
	}

	fmt.Println()
	ui.CurrentTheme.Primary.Println("aka")
	ui.CurrentTheme.Body.Println("Create short commands for launching GUI applications")
	ui.Divider()

	fmt.Println()
	ui.CurrentTheme.Primary.Println("USAGE")
	fmt.Print("  ")
	ui.CurrentTheme.Body.Printf("aka [command]\n")

	if cmd.HasAvailableSubCommands() {
		fmt.Println()
		ui.CurrentTheme.Primary.Println("COMMANDS")

		commands := cmd.Commands()
		for _, c := range commands {
			if !c.IsAvailableCommand() || c.Name() == "help" {
				continue
			}
			ui.CurrentTheme.Accent.Printf("  %-12s", c.Name())
			ui.CurrentTheme.Body.Printf("%s\n", c.Short)
		}
	}

	if cmd.HasAvailableLocalFlags() {
		fmt.Println()
		ui.CurrentTheme.Primary.Println("FLAGS")
		ui.CurrentTheme.Body.Print("  -h, --help       Show this help message\n")
		ui.CurrentTheme.Body.Print("  -v, --version    Show version information\n")
	}

	fmt.Println()
	ui.CurrentTheme.Primary.Println("EXAMPLES")
	ui.PrintCommand("aka add safari Safari")
	ui.PrintCommand("aka add code \"Visual Studio Code\"")
	ui.PrintCommand("aka list")
	ui.PrintCommand("aka rename safari web")
	ui.PrintCommand("aka remove web")

	fmt.Println()
	ui.CurrentTheme.Muted.Printf("Run 'aka [command] --help' for more information about a command.\n")
	fmt.Println()
}

// subcommandHelp displays themed help for subcommands
func subcommandHelp(cmd *cobra.Command) {
	fmt.Println()
	ui.CurrentTheme.Primary.Println(cmd.Name())
	if cmd.Short != "" {
		ui.CurrentTheme.Body.Println(cmd.Short)
	}
	ui.Divider()

	if cmd.Long != "" {
		fmt.Println()
		ui.CurrentTheme.Body.Println(cmd.Long)
	}

	fmt.Println()
	ui.CurrentTheme.Primary.Println("USAGE")
	fmt.Print("  ")
	ui.CurrentTheme.Body.Println(cmd.UseLine())

	if cmd.HasAvailableLocalFlags() {
		fmt.Println()
		ui.CurrentTheme.Primary.Println("FLAGS")
		fmt.Print(cmd.LocalFlags().FlagUsages())
	}

	fmt.Println()
}

func getUsageTemplate() string {
	return `
` + ui.CurrentTheme.Primary.Sprint("{{.Name}}") + `
{{if .Short}}{{.Short}}
` + ui.CurrentTheme.Border.Sprint("──────────────────────────────────────────────────") + `
{{end}}{{if .Long}}
{{.Long | trimTrailingWhitespaces}}
{{end}}
` + ui.CurrentTheme.Primary.Sprint("USAGE") + `
  {{.UseLine}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}
{{if gt (len .Aliases) 0}}

` + ui.CurrentTheme.Primary.Sprint("ALIASES") + `
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

` + ui.CurrentTheme.Primary.Sprint("EXAMPLES") + `
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

` + ui.CurrentTheme.Primary.Sprint("COMMANDS") + `{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  ` + ui.CurrentTheme.Accent.Sprint("{{rpad .Name .NamePadding}}") + ` {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

` + ui.CurrentTheme.Primary.Sprint("{{.Title}}") + `{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  ` + ui.CurrentTheme.Accent.Sprint("{{rpad .Name .NamePadding}}") + ` {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

` + ui.CurrentTheme.Primary.Sprint("ADDITIONAL COMMANDS") + `{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  ` + ui.CurrentTheme.Accent.Sprint("{{rpad .Name .NamePadding}}") + ` {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

` + ui.CurrentTheme.Primary.Sprint("FLAGS") + `
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

` + ui.CurrentTheme.Primary.Sprint("GLOBAL FLAGS") + `
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

` + ui.CurrentTheme.Primary.Sprint("ADDITIONAL HELP TOPICS") + `{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

` + ui.CurrentTheme.Muted.Sprint("Use \"{{.CommandPath}} [command] --help\" for more information about a command.") + `{{end}}
`
}
