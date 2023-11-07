package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	utils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
)

type Logger struct {
	*logrus.Logger
	ctx context.Context
}

func NewLogger(ctx context.Context) *Logger {
	logFolderPath := "logs"
	l := logrus.New()

	l.SetReportCaller(true)

	l.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	err := os.MkdirAll(logFolderPath, 0755)
	if err != nil {
		panic(err)
	}

	lumber := &lumberjack.Logger{
		Filename: logFolderPath + "/server.log",
		MaxSize:  30,
		MaxAge:   2,
		Compress: false,
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.SetLevel(logrus.InfoLevel)

	l.AddHook(&writeHook{
		Writer:    []io.Writer{allFile, lumber, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	return &Logger{
		Logger: l,
		ctx:    ctx,
	}
}

// ------------------ Multi Writer ---------------------
type writeHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writeHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	id := utils.GetReqID(entry.Context)
	entry.Data["Request"] = id

	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writeHook) Levels() []logrus.Level {
	return hook.LogLevels
}
