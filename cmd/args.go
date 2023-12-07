package cmd

import (
	"errors"
	"github.com/devdimensionlab/plybuild/pkg/bitbucket"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/context"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/http"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/shell"
	"github.com/devdimensionlab/plybuild/pkg/spring"
	"github.com/devdimensionlab/plybuild/pkg/template"
	"github.com/spf13/cobra"
	"os"
)

var ctx context.Context

func InitGlobals(cmd *cobra.Command) error {
	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}
	if json {
		logger.SetFieldLogger()
		logger.SetJsonLogging()
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	if debug {
		debugLogger := logger.DebugLogger()
		log = debugLogger
		bitbucket.SetLogger(debugLogger)
		config.SetLogger(debugLogger)
		context.SetLogger(debugLogger)
		file.SetLogger(debugLogger)
		http.SetLogger(debugLogger)
		maven.SetLogger(debugLogger)
		shell.SetLogger(debugLogger)
		spring.SetLogger(debugLogger)
		template.SetLogger(debugLogger)
	}

	stealth, _ := cmd.Flags().GetBool("stealth")
	if stealth {
		logger.SetFieldLogger()
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
