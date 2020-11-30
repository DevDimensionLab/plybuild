package shell

import (
	"github.com/co-pilot-cli/co-pilot/pkg/logger"
	"github.com/sirupsen/logrus"
)

var log = logger.Context()

func SetLogger(logger logrus.FieldLogger) {
	log = logger
}
