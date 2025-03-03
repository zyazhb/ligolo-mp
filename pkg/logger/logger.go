package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
	"time"
)

type LogHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	var logLine strings.Builder

	timeStr := r.Time.Format(time.RFC3339)
	logLine.WriteString(fmt.Sprintf("[#A9A9A9][%s][-]", timeStr))
	logLine.WriteString(" ")

	switch r.Level {
	case slog.LevelDebug:
		logLine.WriteString(fmt.Sprintf("[#FF00FF]%s[-]", r.Level))
	case slog.LevelInfo:
		logLine.WriteString(fmt.Sprintf("[#7CFC00]%s[-]", r.Level))
	case slog.LevelWarn:
		logLine.WriteString(fmt.Sprintf("[#FFFF00]%s[-]", r.Level))
	case slog.LevelError:
		logLine.WriteString(fmt.Sprintf("[#880808]%s[-]", r.Level))
	}
	logLine.WriteString(" ")

	logLine.WriteString(fmt.Sprintf("[#F9F6EE]%s[-]", r.Message))
	logLine.WriteString(" ")

	r.Attrs(func(a slog.Attr) bool {
		logLine.WriteString(fmt.Sprintf("[#87CEEB]%s[-]=[#C2B280]%s[-]", a.Key, a.Value))
		logLine.WriteString(" ")
		return true
	})

	h.l.Println(logLine.String())

	return nil
}

func NewLogHandler(out io.Writer, opts *slog.HandlerOptions) *LogHandler {
	h := &LogHandler{
		Handler: slog.NewTextHandler(out, opts),
		l:       log.New(out, "", 0),
	}

	return h
}
