package sqlconnect

import (
	"database/sql"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToPgDB(dsn string, log *slog.Logger) (*gorm.DB, error) {
	log.Info("connecting to PostgreSQL database")

	gormConfig := &gorm.Config{
		// Caches prepared statements per-connection for reuse; shows as brief "idle in transaction" in pg_stat_activity on new connections — expected behaviour.
		PrepareStmt: true,
		// GORM wraps every write in an implicit transaction by default; disabled because withTx manages transactions explicitly to avoid double-wrapping.
		SkipDefaultTransaction: true,
		Logger:                 newSlogGORMLogger(log, logger.Warn, 200*time.Millisecond),
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		return nil, err
	}

	return gormDB, nil
}

func withTx(db *gorm.DB, fn func(tx *gorm.DB) error, opts *sql.TxOptions) error {
	err := db.Transaction(fn, opts)
	if err != nil {
		slog.Error("Transaction failed", "error", err)
		return err
	}
	return nil

	// tx := db.Begin()
	// if tx.Error != nil {
	// 	return tx.Error
	// }

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	// if err := fn(tx); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// if err := tx.Commit().Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

}
