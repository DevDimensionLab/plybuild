package cmd

import (
	"co-pilot/pkg/service"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ctx service.Context

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
