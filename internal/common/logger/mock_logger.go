package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type CustomLoggerMock struct {
	mock.Mock
}

func (m *CustomLoggerMock) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *CustomLoggerMock) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *CustomLoggerMock) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *CustomLoggerMock) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *CustomLoggerMock) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *CustomLoggerMock) Panic(args ...interface{}) {
	m.Called(args...)
}

func (m *CustomLoggerMock) Debugf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *CustomLoggerMock) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *CustomLoggerMock) Warnf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *CustomLoggerMock) Errorf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *CustomLoggerMock) Fatalf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *CustomLoggerMock) Panicf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *CustomLoggerMock) SetLevel(level logrus.Level) {
	m.Called(level)
}
