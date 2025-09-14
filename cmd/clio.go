package main

import (
	"fmt"
	"os"

	"github.com/muixstudio/clio/cmd/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clio",
	Short: "clio short",
	Long:  `clio long`,
}

func init() {
	rootCmd.AddCommand(service.Cmd)
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
