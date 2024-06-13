package cmd

import (
	log "log/slog"
	"os"

	"github.com/omran95/chatroom/internal/wire"
	"github.com/spf13/cobra"
)

var subscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "Subscriber Service",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeSubscriberServer("subscriber")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	appCmd.AddCommand(subscriberCmd)
}
