package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func NewLogger(config Config) (*Logger, error) {

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("failed to create folder: %v", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal level: %v", err)
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				return slog.String(slog.TimeKey, t.Format("2006-01-02T15:04:05.000000"))
			}

			return a
		},
	}

	handlers := slog.NewMultiHandler(
		slog.NewJSONHandler(logFile, opts),
		slog.NewTextHandler(os.Stdout, opts),
	)

	logger := slog.New(handlers)

	return &Logger{logger, logFile}, nil
}

func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value("logger").(*Logger)
	if !ok {
		panic("logger not found in context")
	}

	return logger
}

func MustNewLogger(config Config) *Logger {
	logger, err := NewLogger(config)

	if err != nil {
		err := fmt.Errorf("failed to create logger: %v", err)
		panic(err)
	}

	return logger
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Printf("failed to close logger: %v", err)
	}
}

func (l *Logger) With(fields ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}
