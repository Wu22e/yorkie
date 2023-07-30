package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yorkie-team/yorkie/cmd/yorkie/config"
)

var (
	flagForce bool
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out from the Yorkie server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read in config: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := config.Load()
			if err != nil {
				return err
			}
			rpcAddr := viper.GetString("rpcAddr")
			if flagForce {
				return config.Delete()
			}
			if rpcAddr == "" {
				return errors.New("you must specify the server address to log out")
			}
			authToken, ok := conf.Auths[rpcAddr]
			if !ok || authToken == "" {
				return fmt.Errorf("you are not logged in to %s", rpcAddr)
			}
			if len(conf.Auths) <= 1 {
				return config.Delete()
			}
			delete(conf.Auths, rpcAddr)
			conf.IsInsecure = false
			conf.RPCAddr = ""
			return config.Save(conf)
		},
	}
}

func init() {
	cmd := newLogoutCmd()
	cmd.Flags().BoolVar(
		&flagForce,
		"force",
		false,
		"force log out from all servers",
	)
	rootCmd.AddCommand(cmd)
}
