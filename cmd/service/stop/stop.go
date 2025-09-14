package stop

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "stop",
	Short: "stop long",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop service")
	},
}
