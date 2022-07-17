package cmd

import (
	"zerologix-homework/src/server"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "啟動伺服器",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
