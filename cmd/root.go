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
		taskdata, taskrc := task.GetEnv()
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

		// Read projects
		projects, err := task.ReadProjects(taskdata + "/projects.data")
		if err != nil {
			fmt.Println("Error reading projects:", err)
			os.Exit(1)
		}
		// TODO: Parse `task _projects` to get additional projects if present

		// Read contexts
		contexts, err := task.ReadContexts(taskdata + "/contexts.data")
		if err != nil {
			fmt.Println("Error reading contexts:", err)
			os.Exit(1)
		}
		// TODO: Parse taskrc to get additional contexts if present

		// TODO: Assuming projects and contexts are corresponding ( p[i] => c[i] )

		// Start the actual program
		if _, err := tea.NewProgram(ui.New(projects, contexts)).Run(); err != nil {
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
