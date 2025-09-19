package main

import (
	"fmt"
	"os"

	"github.com/muixstudio/clio/tools/gen-model/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gen-model",
	Short: "clio short",
	Long:  `clio long`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.GenModel(cmd.Flag("sql_file").Value.String())
	},
}

func init() {
	rootCmd.Flags().StringP("model_name", "n", "", "model name")
	rootCmd.Flags().StringP("sql_file", "", "", "sql file")
	rootCmd.Flags().StringP("table_prefix", "", "", "table prefix")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
