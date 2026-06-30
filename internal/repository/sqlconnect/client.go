package sqlconnect

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToPgDB(dsn string, log *slog.Logger) (*gorm.DB, error) {
	log.Info("connecting to PostgreSQL database")

	const maxRetries = 5
	baseDelay := 5 * time.Second
	maxDelay := 60 * time.Second

	gormConfig := &gorm.Config{
		// Caches prepared statements per-connection for reuse; shows as brief "idle in transaction" in pg_stat_activity on new connections — expected behaviour.
		PrepareStmt: true,
		// GORM wraps every write in an implicit transaction by default; disabled because withTx manages transactions explicitly to avoid double-wrapping.
		SkipDefaultTransaction: true,
		Logger:                 newSlogGORMLogger(log, logger.Warn, 200*time.Millisecond),
	}

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, openErr := gorm.Open(postgres.Open(dsn), gormConfig)
		if openErr != nil {
			lastErr = openErr
			log.Warn("Database connection attempt failed", "attempt", attempt, "error", openErr)
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				log.Error("Failed to get sql.DB from gorm.DB", "error", err)
				return nil, err
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			pingErr := sqlDB.PingContext(ctx)
			cancel()
			if pingErr == nil {
				log.Info("Successfully connected to the database")
				return db, nil
			}
			lastErr = pingErr
			log.Error("Failed to ping database", "error", pingErr)
		}

		if attempt == maxRetries {
			break
		}
		delay := min(baseDelay*time.Duration(1<<uint(attempt-1)), maxDelay)
		log.Info("Retrying database connection", "attempt", attempt, "delay", delay.String())
		time.Sleep(delay)
	}

	log.Error("Max retries reached. Could not connect to the database.")
	return nil, lastErr
}

func withTx(db *gorm.DB, fn func(tx *gorm.DB) error, opts *sql.TxOptions) error {
	err := db.Transaction(fn, opts)
	if err != nil {
		slog.Error("Transaction failed", "error", err)
		return err
	}
	return nil
}
