package cmd

import (
	"fmt"
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
		cfg, err := actions.ParseRc(taskrc)

		contexts := actions.RcExtract("context", cfg)
		for k, v := range contexts {
			fmt.Println(k, "---", v)
		}

		reports := actions.RcExtract("report", cfg)
		for k, v := range reports {
			fmt.Println(k, "---", v)
		}

		urgency := actions.RcExtract("urgency", cfg)
		for k, v := range urgency {
			fmt.Println(k, "---", v)
		}

		if err != nil {
			log.Fatalln(err)
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
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
