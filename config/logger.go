package config

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var Logger AmarthaLogger

type AmarthaLogger struct {
	zerolog.Logger
}

func NewLogger() AmarthaLogger {
	l := zerolog.New(os.Stdout)
	return AmarthaLogger{
		Logger: l,
	}
}

func (l *AmarthaLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.Error().Msgf(format, v...)
}

func (l *AmarthaLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.Info().Msgf(format, v...)
}

func (l *AmarthaLogger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.Warn().Msgf(format, v...)
}
