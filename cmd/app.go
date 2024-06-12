package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appCmd = &cobra.Command{
	Use:   "root",
	Short: "chat app executable",
}

func Execute() {
	cobra.CheckErr(appCmd.Execute())
}

var configFile string

func init() {
	cobra.OnInitialize(initConfig)

	appCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")
}

func initConfig() {
	viper.SetConfigType("yaml")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
