package logs

import (
	"context"
	"fmt"
	"log/slog"
)

type Logger interface {
	Info(ctx context.Context, message string)
	Infof(ctx context.Context, format string, args ...any)
	Warn(ctx context.Context, message string)
	Warnf(ctx context.Context, format string, args ...any)
	Error(ctx context.Context, message string)
	Errorf(ctx context.Context, format string, args ...any)
}

type logger struct {
	LoggerOpts

	slog *slog.Logger
}

func NewLogger(opts LoggerOpts) Logger {
	var handler slog.Handler

	switch opts.Style {
	case StyleJSON:
		handler = slog.NewJSONHandler(opts.Out, &slog.HandlerOptions{
			Level: opts.Level,
		})
	default:
		handler = slog.NewTextHandler(opts.Out, &slog.HandlerOptions{
			Level: opts.Level,
		})
	}

	return &logger{
		LoggerOpts: opts,

		slog: slog.New(handler),
	}
}

func (l *logger) getAnnotations(ctx context.Context) []any {
	annotations := GetAnnotations(ctx)

	loggables := []any{}
	for key, value := range annotations {
		loggables = append(loggables, key, value)
	}

	return loggables
}

func (l *logger) Info(ctx context.Context, message string) {
	l.slog.InfoContext(ctx, message, l.getAnnotations(ctx)...)
}

func (l *logger) Infof(ctx context.Context, format string, args ...any) {
	l.Info(ctx, fmt.Sprintf(format, args...))
}

func (l *logger) Warn(ctx context.Context, message string) {
	l.slog.WarnContext(ctx, message, l.getAnnotations(ctx)...)
}

func (l *logger) Warnf(ctx context.Context, format string, args ...any) {
	l.Warn(ctx, fmt.Sprintf(format, args...))
}

func (l *logger) Error(ctx context.Context, message string) {
	l.slog.ErrorContext(ctx, message, l.getAnnotations(ctx)...)
}

func (l *logger) Errorf(ctx context.Context, format string, args ...any) {
	l.Error(ctx, fmt.Sprintf(format, args...))
}
