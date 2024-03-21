package cmd

import (
	"bufio"
	"fmt"
	// "log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	ui "github.com/manyids2/tasktea/components/navbar"
	"github.com/spf13/cobra"
)

var history_path string
var config_path string

func LinesFromFile(logfile string) []string {
	lines := []string{}
	f, err := os.OpenFile(logfile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return lines
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return lines
	}
	return lines
}

// taskTeaCmd Launches application
var taskTeaCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "Taskwarrior TUI",
	Long: `Taskwarrior TUI.
Set filters to view tasks as trees with dependencies.`,
	Run: func(cmd *cobra.Command, filters []string) {

		// // Load history
		// history := LinesFromFile(history_path)
		//
		// // Prep to write to history
		// f, err := os.OpenFile(history_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		// if err != nil {
		// 	log.Fatalf("error opening file: %v", err)
		// }
		// defer f.Close()
		// log.SetOutput(f)
		// log.SetFlags(0)

		// if _, err := tea.NewProgram(ui.NewModel(filters, history)).Run(); err != nil {
		// 	fmt.Printf("Could not start program :(\n%v\n", err)
		// 	os.Exit(1)
		// }

		if _, err := tea.NewProgram(ui.NewModel("context", filters)).Run(); err != nil {
			fmt.Printf("Could not start program :(\n%v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := taskTeaCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Icons, taskrc etc.
	taskTeaCmd.Flags().StringVarP(&config_path, "config", "C",
		"/home/x/.config/tasktea/config.yaml", "Config file")

	// Filter history
	taskTeaCmd.Flags().StringVarP(&history_path, "history", "H",
		"/home/x/.config/tasktea/history", "History file")
}
