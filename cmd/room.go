package cmd

import (
	log "log/slog"
	"os"

	"github.com/omran95/chat-app/internal/wire"
	"github.com/spf13/cobra"
)

var roomCmd = &cobra.Command{
	Use:   "room",
	Short: "Room Service",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeRoomServer("room")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	appCmd.AddCommand(roomCmd)
}
