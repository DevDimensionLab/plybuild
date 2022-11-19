package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var fieldLogger = false
var collector = &Collector{}
var log = logrus.New()

func init() {
	log.AddHook(collector)
}

func DebugLogger() logrus.FieldLogger {
	pc, file, line, ok := runtime.Caller(1)

	var fields logrus.Fields

	if ok {
		fields = logrus.Fields{
			"caller": file,
			"line":   line,
			"func":   runtime.FuncForPC(pc).Name(),
		}
	}

	log.SetLevel(logrus.DebugLevel)
	return log.WithFields(fields)
}

func Context() logrus.FieldLogger {
	return log.WithFields(logrus.Fields{})
}

func ExternalError(err error, msg string) error {
	output := fmt.Sprintf("%v\n\n##### EXTERNAL ERROR MESSAGE #####\n\n", err)

	for _, line := range strings.Split(msg, "\n") {
		output += fmt.Sprintf("  %s\n", line)
	}

	output = fmt.Sprintf("%s##### ENDS EXTERNAL ERROR MESSAGE #####\n\n", output)

	return errors.New(output)
}

func SetJsonLogging() {
	log.SetFormatter(&logrus.JSONFormatter{})
}

func SetFieldLogger() {
	fieldLogger = true
}

func IsFieldLogger() bool {
	return fieldLogger
}

func StdOut() *os.File {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		return os.Stdout
	}
	return nil
}

func LogEntries() []*logrus.Entry {
	return collector.entries
}
