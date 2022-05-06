package spring

import (
	"github.com/devdimensionlab/co-pilot/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log = logger.Context()

func SetLogger(logger logrus.FieldLogger) {
	log = logger
}
