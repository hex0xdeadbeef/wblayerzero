package postgres

import (
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

	tx, err := db.db.Beginx()
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
func (db *Database) GetAllOrders() (res []entities.Order, err error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("starting a new tx: %w", err)
	}

	const (
		query = `
		SELECT
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
			o.customer_id, o.delivery_service, o.date_created, o.shardkey, o.sm_id, o.oof_shard,
			i.order_uid as item_order_uid, i.track_number as item_track_number, i.status as item_status,
			i.chrt_id as item_chrt_id, i.nm_id as item_nm_id, i.rid as item_rid,
			i.brand as item_brand, i.name as item_name, i.size as item_size,
			i.price as item_price, i.sale as item_sale, i.total_price as item_total_price,
			p.order_uid as payment_order_uid, p.transaction_id, p.request_id, p.bank, p.currency,
			p.provider, p.payment_dt, p.amount, p.delivery_cost, p.goods_total, p.custom_fee,
			d.order_uid as delivery_order_uid, d.zip, d.name, d.phone, d.email,
			d.city, d.address, d.region
		FROM orders o, items i, payments p, deliveries d
		WHERE o.order_uid = i.order_uid
		  AND o.order_uid = p.order_uid
		  AND o.order_uid = d.order_uid;
	`
		startSize = 1 << 8
	)

	rows, err := tx.Queryx(query)
	if err != nil {
		return nil, wrapTxError(tx, "getting all orders", err)
	}
	defer func() {
		rowsCloseErr := rows.Close()
		if rowsCloseErr == nil {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; closing rows: %w", err, rowsCloseErr)
		}
		err = rowsCloseErr
	}()

	var (
		ordersMap = make(map[string]entities.Order, startSize)

		curOrder    entities.Order
		curItem     entities.Item
		curPayment  entities.Payment
		curDelivery entities.Delivery
	)
	for rows.Next() {
		if err = rows.Scan(
			&curOrder.UID, &curOrder.TrackNumber, &curOrder.Entry, &curOrder.Locale, &curOrder.InternalSignature, &curOrder.CustomerID, &curOrder.DeliveryService,
			&curOrder.DateCreated, &curOrder.ShardKey, &curOrder.SmID, &curOrder.OffShard,

			&curItem.OrderUID, &curItem.TrackNumber, &curItem.Status, &curItem.ChrtID, &curItem.NmID, &curItem.RID, &curItem.Brand, &curItem.Name, &curItem.Size,
			&curItem.Price, &curItem.Sale, &curItem.TotalPrice,

			&curPayment.OrderUID, &curPayment.TxID, &curPayment.RequestID, &curPayment.Bank, &curPayment.Currency, &curPayment.Provider, &curPayment.PaymentDt,
			&curPayment.Amount, &curPayment.DeliveryCost, &curPayment.GoodsTotal, &curPayment.CustomFee,

			&curDelivery.OrderUID, &curDelivery.Zip, &curDelivery.Name, &curDelivery.Phone, &curDelivery.Email,
			&curDelivery.City, &curDelivery.Address, &curDelivery.Region); err != nil {
			return nil, wrapTxError(tx, "scanning row into struct", err)
		}

		presentedOrder, ok := ordersMap[curOrder.UID]
		if ok {
			presentedOrder.Items = append(presentedOrder.Items, curItem)
			continue
		}

		curOrder.Items, curOrder.Payment, curOrder.Delivery = append(curOrder.Items, curItem), curPayment, curDelivery
		ordersMap[curOrder.UID] = curOrder

	}
	if err := rows.Err(); err != nil {
		return nil, wrapTxError(tx, "after rows scanning: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, wrapTxError(tx, "committing tx:", err)
	}

	res = make([]entities.Order, 0, len(ordersMap))
	for _, v := range ordersMap {
		res = append(res, v)
	}

	return res, nil
}

// prepareAndExecStmt prepares a stmt with the query param and if everyrhing is okay executes it (otherwise returns a wrapped error)
// returns an error if any
func prepareAndExecStmt(tx *sqlx.Tx, query string, args ...any) error {
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
func wrapTxError(tx *sqlx.Tx, subErrMsgF string, subError error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return fmt.Errorf("%s: %w; during rollback: %w", subErrMsgF, subError, rollbackErr)
	}

	return fmt.Errorf(subErrMsgF, subError)
}
