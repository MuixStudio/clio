package service

import (
	"fmt"

	"github.com/muixstudio/clio/cmd/service/restart"
	"github.com/muixstudio/clio/cmd/service/start"
	"github.com/muixstudio/clio/cmd/service/stop"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "service",
	Short: "service long",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service service")
	},
}

func init() {
	Cmd.AddCommand(start.Cmd)
	Cmd.AddCommand(restart.Cmd)
	Cmd.AddCommand(stop.Cmd)
}
