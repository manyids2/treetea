package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manyids2/tasktea/ui"
	"github.com/spf13/cobra"
)

// taskTeaCmd Launches application
var taskTeaCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "Taskwarrior TUI",
	Long: `Taskwarrior TUI.
Set filters to view tasks as trees with dependencies.`,
	Run: func(cmd *cobra.Command, filters []string) {
		if _, err := tea.NewProgram(ui.NewModel(filters)).Run(); err != nil {
			fmt.Printf("Could not start program :(\n%v\n", err)
			os.Exit(1)
		}

		// TODO: Load config

		// config, err := cmd.Flags().GetString("config")
		// if err != nil {
		// 	fmt.Printf("Could not read config argument :(\n%v\n", err)
		// 	os.Exit(1)
		// }
		// fmt.Println("config:", config)
	},
}

func Execute() {
	err := taskTeaCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	taskTeaCmd.Flags().StringP("config", "C", "$XDG_DATA_DIR/tasktea/config.yaml", "Config file")
}
