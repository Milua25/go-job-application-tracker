package sqlconnect

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

type slogGORMLogger struct {
	log           *slog.Logger
	level         logger.LogLevel
	slowThreshold time.Duration
}

// newSlogGORMLogger creates a new slogGORMLogger that implements the gorm/logger.Interface.
func newSlogGORMLogger(log *slog.Logger, level logger.LogLevel, slowThreshold time.Duration) logger.Interface {
	return &slogGORMLogger{
		log:           log,
		level:         level,
		slowThreshold: slowThreshold,
	}
}

func (l *slogGORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &slogGORMLogger{log: l.log, level: level, slowThreshold: l.slowThreshold}
}

func (l *slogGORMLogger) Info(_ context.Context, msg string, args ...any) {
	if l.level >= logger.Info {
		l.log.Info(msg, args...)
	}
}

func (l *slogGORMLogger) Warn(_ context.Context, msg string, args ...any) {
	if l.level >= logger.Warn {
		l.log.Warn(msg, args...)
	}
}

func (l *slogGORMLogger) Error(_ context.Context, msg string, args ...any) {
	if l.level >= logger.Error {
		l.log.Error(msg, args...)
	}
}

func (l *slogGORMLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []any{
		slog.Duration("elapsed", elapsed),
		slog.String("sql", sql),
		slog.Int64("rows", rows),
	}

	switch {
	// ErrRecordNotFound is a normal flow return from First/Take; callers handle it, so suppress it here to avoid noise.
	case err != nil && !errors.Is(err, logger.ErrRecordNotFound):
		l.log.Error("gorm query error", append(attrs, slog.String("error", err.Error()))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold:
		l.log.Warn("gorm slow query", attrs...)
	case l.level >= logger.Info:
		l.log.Debug("gorm query", attrs...)
	}
}
