package http

import (
	"github.com/devdimensionlab/ply/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log = logger.Context()

func SetLogger(logger logrus.FieldLogger) {
	log = logger
}
