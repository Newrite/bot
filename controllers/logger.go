package controllers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
)

type myFormatter struct {
	logrus.TextFormatter
}

var logger = &logrus.Logger{}
var once sync.Once

func SingleLog() *logrus.Logger {
	once.Do(func() {
		logger = newLogger()
	})
	return logger
}

func (f *myFormatter) format(entry *logrus.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the logrus levels
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 // gray
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return []byte(fmt.Sprintf("[%s] - \x1b[%dm%s\x1b[0m - %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

func newLogger() *logrus.Logger {
	f, _ := os.OpenFile("logrus.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stderr, f),
		Level: logrus.InfoLevel,
		Formatter: &myFormatter{logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            true,
			DisableLevelTruncation: true,
		},
		},
	}
	return logger
}
