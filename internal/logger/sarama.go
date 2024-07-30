package logger

import (
	"context"
	"fmt"
	"log/slog"
)

type SaramaLogger struct {
	log   *slog.Logger
	level slog.Level
}

func NewSaramaLogger(log *slog.Logger, level slog.Level) *SaramaLogger {
	return &SaramaLogger{log: log.With(slog.String("lib", "sarama")), level: level}
}

func (l *SaramaLogger) Print(v ...interface{}) {
	l.log.Log(context.Background(), l.level, fmt.Sprint(v...))
}

func (l *SaramaLogger) Printf(format string, v ...interface{}) {
	l.log.Log(context.Background(), l.level, fmt.Sprintf(format, v...))
}

func (l *SaramaLogger) Println(v ...interface{}) {
	l.log.Log(context.Background(), l.level, fmt.Sprintln(v...))
}
