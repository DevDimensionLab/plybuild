package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var (
	Info = White
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func Context() logrus.FieldLogger {
	_, _, _, ok := runtime.Caller(1)

	var fields logrus.Fields

	if ok {
		fields = logrus.Fields{
			//"caller": file,
			//"line": line,
			//"func": runtime.FuncForPC(pc).Name(),
		}
	}
	return logrus.WithFields(fields)
}

func ExternalError(err error, msg string) error {
	output := fmt.Sprintf("%v\n\n##### EXTERNAL ERROR MESSAGE #####\n\n", err)

	for _, line := range strings.Split(msg, "\n") {
		output += fmt.Sprintf("  %s\n", line)
	}

	output = fmt.Sprintf("%s##### ENDS EXTERNAL ERROR MESSAGE #####\n\n", output)

	return errors.New(output)
}

func StdOut() *os.File {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		return os.Stdout
	}
	return nil
}
