package spring

import (
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log = logger.Context()

func SetLogger(logger logrus.FieldLogger) {
	log = logger
}
