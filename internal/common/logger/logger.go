package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type CustomLogger struct {
	*logrus.Logger
}

func CreateCustomLogger() *CustomLogger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return &CustomLogger{Logger: logger}
}
