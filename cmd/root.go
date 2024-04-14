package cmd

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manyids2/tasktea/task"
	"github.com/manyids2/tasktea/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "Tools for taskwarrior",
	Long:  `Like a standalone taskwiki`,
	Run: func(cmd *cobra.Command, args []string) {
		_, taskrc := task.GetEnv()
		if !task.FileExists(taskrc) {
			fmt.Println("Initializing taskrc:", taskrc)

			// Initialize program, actually dont need tea, but meh
			cmd := exec.Command("task")
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			// Run command
			cmd.Run()
			if !task.FileExists(taskrc) {
				fmt.Println("Need taskrc to continue.")
				os.Exit(1)
			}
		}

		// Start the actual program
		if _, err := tea.NewProgram(ui.New()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Like nvim -u ...
	rootCmd.PersistentFlags().StringP("config", "u", "$TASKRC/tasktea.yaml", "Config file (default is $TASKRC/tasktea.yaml)")
}
