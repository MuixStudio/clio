package restart

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "restart",
	Short: "restart long",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restart service")
	},
}
