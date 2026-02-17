package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)

	return &LogrusLogger{
		logger: logger,
	}
}

func (ll *LogrusLogger) Info(ctx context.Context, message string) {
	ll.logger.Info(message)
}

func (ll *LogrusLogger) Error(ctx context.Context, message string) {
	ll.logger.Error(message)
}

func (ll *LogrusLogger) Debug(ctx context.Context, message string) {
	ll.logger.Debug(message)
}

func (ll *LogrusLogger) Warn(ctx context.Context, message string) {
	ll.logger.Warn(message)
}

func (ll *LogrusLogger) Fatal(ctx context.Context, message string) {
	ll.logger.Fatal(message)
}
