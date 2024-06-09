-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
	order_uid VARCHAR,
	track_number VARCHAR UNIQUE NOT NULL,

	entry VARCHAR NOT NULL,
	locale VARCHAR NOT NULL,
	internal_signature VARCHAR NOT NULL,

	customer_id VARCHAR NOT NULL,
	delivery_service VARCHAR NOT NULL,

	shardkey INT NOT NULL,
	sm_id INT NOT NULL,
	date_created DATE NOT NULL,
	oof_shard INT NOT NULL,

	PRIMARY KEY (order_uid)
);

CREATE TABLE IF NOT EXISTS items (
	order_uid VARCHAR PRIMARY KEY,

	track_number VARCHAR NOT NULL,

	status INT NOT NULL,

	chrt_id INT NOT NULL,
	nm_id INT NOT NULL,
    rid VARCHAR NOT NULL,

	brand VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
	size VARCHAR,

	price NUMERIC(30, 15) NOT NULL,
	sale NUMERIC(30, 15) NOT NULL,
	total_price NUMERIC(30, 15) NOT NULL,

    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
	order_uid VARCHAR PRIMARY KEY,

	transaction_id VARCHAR UNIQUE NOT NULL,
	request_id VARCHAR UNIQUE NOT NULL,

	bank VARCHAR NOT NULL,
	currency VARCHAR NOT NULL,
	provider VARCHAR NOT NULL,

	payment_dt TIMESTAMP NOT NULL,

	amount NUMERIC(30, 15) NOT NULL,
	delivery_cost NUMERIC(30, 15) NOT NULL,
	goods_total NUMERIC(30, 15) NOT NULL,
	custom_fee NUMERIC(30, 15) NOT NULL,

	FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS deliveries (
	order_uid VARCHAR PRIMARY KEY,

	name VARCHAR NOT NULL,
	phone VARCHAR NOT NULL,
	email VARCHAR NOT NULL,

	zip VARCHAR NOT NULL,
	city VARCHAR NOT NULL,
	address VARCHAR NOT NULL,
	region VARCHAR NOT NULL,

	FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
