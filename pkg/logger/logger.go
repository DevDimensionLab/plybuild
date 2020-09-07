package logger

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

func Context() *logrus.Entry {
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
