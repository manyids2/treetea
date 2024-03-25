package cmd

import (
	"fmt"
	"log"
	"os"

	tw "github.com/manyids2/tasktea/task"
	"github.com/manyids2/tasktea/task/actions"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, taskrc := actions.Paths()
		rc, err := tw.LoadTaskRC(taskrc)
		fmt.Println(rc)
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
