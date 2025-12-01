package cmd

import (
	"log/slog"
	"os"

	"github.com/Wangch29/ikun-messenger/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ikun",
	Short: "Ikun Messenger is a distributed IM system, based on Raft",
	Long:  "Ikun Messenger is a distributed IM system, based on Raft",
}

func Execute() {
	initConfig()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file:", "file", viper.ConfigFileUsed())
		viper.Unmarshal(&config.Global)
	}
}
