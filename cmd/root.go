package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manyids2/tasktea/components/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "Minimal UI for Taskwarrior",
	Long:  `Minimal UI for Taskwarrior based on taskwiki, implemented using bubbletea.`,
	Run: func(cmd *cobra.Command, args []string) {
		tt := app.Model{}
		tt.LoadRc()
		if _, err := tea.NewProgram(tt).Run(); err != nil {
			fmt.Printf("Could not start program :(\n%v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func init() {}
