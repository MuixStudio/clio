package start

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "strat",
	Short: "start long",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start service")
	},
}
