package cmd

import (
	"log"
	"os"

	"github.com/manyids2/tasktea/task/actions"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, taskrc := actions.Paths()
		_, err := actions.ParseRc(taskrc)

		if err != nil {
			log.Fatalln(err)
		}
	},
}

func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func init() {}
