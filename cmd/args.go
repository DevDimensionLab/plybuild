package cmd

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var ctx config.Context

func EnableDebug(cmd *cobra.Command) error {
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	if debug {
		fmt.Println("== debug mode enabled ==")
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func OkHelp(cmd *cobra.Command, depend func() bool) error {
	if !cmd.Flags().HasFlags() || !depend() {
		_ = cmd.Help()
		os.Exit(0)
	}
	return nil
}

func SyncActiveProfileCloudConfig() error {
	if ctx.ForceCloudSync {
		if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
			return errors.New("failed to sync cloud config: " + err.Error())
		}
	}
	return nil
}
