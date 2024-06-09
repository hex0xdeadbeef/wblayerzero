package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"wblayerzero/internal/config"
	"wblayerzero/internal/entities"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type (
	Database struct {
		db     *sqlx.DB
		logger *log.Logger
	}
)

// New returns an object representing an abstraction over a DB conn and an errors if any
func New(logger *log.Logger) (*Database, error) {
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

	return &Database{db: db, logger: logger}, nil
}

// Close closes the DB instance and returns error if any
func (db *Database) Close() error {
	if err := db.db.Close(); err != nil {
		return fmt.Errorf("closing db: %w", err)
	}

	return nil
}

// InsertOrder inserts a given order into db and returns an error if any
func (db *Database) InsertOrder(order entities.Order) error {
	var (
		queries = map[string]string{
			"order": `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, date_created, shardkey, sm_id, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,

			"item": `INSERT INTO items (order_uid, track_number, status, chrt_id, nm_id, rid, brand, name, size, price, sale, total_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`,

			"payment": `INSERT INTO payments (order_uid, transaction_id, request_id, bank, currency, provider, payment_dt, amount, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,

			"delivery": `INSERT INTO deliveries (order_uid, name, phone, email, zip, city, address, region) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		}
	)

	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("starting a new tx: %w", err)
	}

	if err := prepareAndExecStmt(tx, queries["order"],
		order.UID,
		order.TrackNumber,

		order.Entry,
		order.Locale,
		order.InternalSignature,

		order.CustomerID,
		order.DeliveryService,

		order.DateCreated,

		order.ShardKey,
		order.SmID,
		order.OffShard); err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	for _, item := range order.Items {
		if err := prepareAndExecStmt(tx, queries["item"],
			item.OrderUID,

			item.TrackNumber,

			item.Status,

			item.ChrtID,
			item.NmID,
			item.RID,

			item.Brand,
			item.Name,
			item.Size,

			item.Price,
			item.Sale,
			item.TotalPrice); err != nil {
			return fmt.Errorf("executing query: %w", err)
		}
	}

	if err := prepareAndExecStmt(tx, queries["payment"],
		order.UID,

		order.TxID,
		order.RequestID,

		order.Bank,
		order.Currency,
		order.Provider,

		order.PaymentDt,

		order.Amount,
		order.DeliveryCost,
		order.GoodsTotal,
		order.CustomFee); err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if err := prepareAndExecStmt(tx, queries["delivery"],
		order.UID,

		order.Name,
		order.Phone,
		order.Email,

		order.Zip,
		order.City,
		order.Address,
		order.Region); err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return wrapTxError(tx, "committing tx:", err)
	}

	return nil
}

// GetAllOrders finds all the orders in DB and returns it if any orders found (otherwise the firs param will be nil) and error if any
func (db *Database) GetAllOrders() ([]entities.Order, error) {

	return nil, nil
}

// prepareAndExecStmt prepares a stmt with the query param and if everyrhing is okay executes it (otherwise returns a wrapped error)
// returns an error if any
func prepareAndExecStmt(tx *sql.Tx, query string, args ...any) error {
	stmt, err := tx.Prepare(query)
	if err != nil {
		return wrapTxError(tx, "preparing stmt:", err)
	}

	if _, err = stmt.Exec(args...); err != nil {
		return wrapTxError(tx, "executing stmt:", err)
	}

	return nil
}

// wrapTxError tries to close the tx and if it's been closed successfully returns only the formatted subError with subErrMsgF,
// otherwise returns a chain of error where both an closing error tx and subError are presented
func wrapTxError(tx *sql.Tx, subErrMsgF string, subError error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return fmt.Errorf("%s: %w; during rollback: %w", subErrMsgF, subError, rollbackErr)
	}

	return fmt.Errorf(subErrMsgF, subError)
}
