package cmd

import (
	"fmt"
	"os"

	"github.com/manyids2/tasktea/components/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasktea",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tt := app.Model{}
		tt.LoadRc()
		fmt.Println(tt.Context, tt.Filters.Read, tt.Projects, tt.Active)
	},
}

func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func init() {}
