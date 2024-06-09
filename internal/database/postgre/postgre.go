package postgre

import (
	"fmt"
	"wblayerzero/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/pressly/goose/v3"
)

type (
	Database struct {
		*sqlx.DB
	}
)

// Close closes the DB instance and returns error if any
func (db *Database) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("closing db: %w", err)
	}

	return nil
}

// New returns an object representing an abstraction over a DB conn and an errors if any
func New() (*Database, error) {
	const (
		driverName     = "pgx"
		migrationsPath = "../../internal/database/migrations"
	)

	db, err := sqlx.Connect(driverName, config.Cfg.GenURI())
	if err != nil {
		return nil, fmt.Errorf("opening and verifying a new conn: %w", err)
	}

	if err := goose.Up(db.DB, migrationsPath); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("applying migrations: %w; closing db: %w", err, closeErr)
		}

		return nil, fmt.Errorf("applying migrations: %w", err)
	}

	return &Database{DB: db}, nil
}
