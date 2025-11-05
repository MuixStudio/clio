package main

import (
	"fmt"
	"os"

	"github.com/muixstudio/clio/services/web/cmd/run"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "web",
	Short: "web short",
	Long:  `web long`,
}

func init() {
	rootCmd.AddCommand(run.Cmd)
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
