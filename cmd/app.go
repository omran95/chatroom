package cmd

import (
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "root",
	Short: "chat app executable",
}

func Execute() {
	cobra.CheckErr(appCmd.Execute())
}
