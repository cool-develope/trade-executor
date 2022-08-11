package storage

import (
	"database/sql"
	"errors"
	"os"

	"github.com/cool-develope/trade-executor/internal/orderctrl/pb"
	_ "github.com/mattn/go-sqlite3" //nolint
)

// Sqlite3Storage is the sqlite3 storage.
type Sqlite3Storage struct {
	db *sql.DB
}

const dbPath = "./storage.db"

// NewSqlite3Storage returns new Sqlite3Storage.
func NewSqlite3Storage() (*Sqlite3Storage, error) {
	isExist := true
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		isExist = false
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	s := &Sqlite3Storage{
		db: db,
	}

	if !isExist {
		err = s.InitDB()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// InitDB inits the db with the schema.
func (s *Sqlite3Storage) InitDB() error {
	_, err := s.db.Exec(`
		DROP TABLE IF EXISTS order_book;
		CREATE TABLE order_book (
			id INTEGER PRIMARY KEY,
			symbol TEXT,
			ask_price REAL,
			ask_qty REAL,
			bid_price REAL,
			bid_qty REAL
		);

		DROP TABLE IF EXISTS applied_order;
		CREATE TABLE applied_order (
			id INTEGER PRIMARY KEY,
			symbol TEXT,
			dir TEXT,
			price REAL,
			amount REAL
		);
		
		DROP TABLE IF EXISTS executed_order;
		CREATE TABLE executed_order (
			price REAL,
			amount REAL,
			order_id INTEGER,
			book_id INTEGER,
			CONSTRAINT fk_order_id
				FOREIGN KEY (order_id)
				REFERENCES applied_order(id),
			CONSTRAINT fk_book_id
				FOREIGN KEY (book_id)
				REFERENCES order_book(id)
		);
		DELETE FROM executed_order;
	`)
	return err
}

// SetOrderBook adds new order book data.
func (s *Sqlite3Storage) SetOrderBook(ob *pb.OrderBook) error {
	_, err := s.db.Exec(`
		INSERT INTO order_book(id, symbol, ask_price, ask_qty, bid_price, bid_qty)
		VALUES(?, ?, ?, ?, ?, ?);`,
		ob.OrderBookId,
		ob.Symbol,
		ob.AskPrice,
		ob.AskQty,
		ob.BidPrice,
		ob.BidQty,
	)

	return err
}

// GetLastOrderBook gets the last order book data.
func (s *Sqlite3Storage) GetLastOrderBook() (*pb.OrderBook, error) {
	var ob pb.OrderBook
	err := s.db.QueryRow(`
		SELECT id, symbol, ask_price, ask_qty, bid_price, bid_qty
		FROM order_book ORDER BY id DESC LIMIT 1;
	`).Scan(&ob.OrderBookId, &ob.Symbol, &ob.AskPrice, &ob.AskQty, &ob.BidPrice, &ob.BidQty)

	return &ob, err
}

// SetAppliedOrder stores the applied order data.
func (s *Sqlite3Storage) SetAppliedOrder(ao *pb.Order) (uint64, error) {
	var id uint64
	tx, err := s.db.Begin()
	if err != nil {
		return id, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO applied_order (symbol, dir, price, amount)
		VALUES (?, ?, ?, ?)
		RETURNING id`)
	if err != nil {
		return id, err
	}
	defer stmt.Close() //nolint

	err = stmt.QueryRow(ao.Symbol, ao.OrderType, ao.Price, ao.Qty).Scan(&id)
	if err != nil {
		return id, err
	}

	err = tx.Commit()
	return id, err
}

// SetExecutedOrder stores the partial executed order.
func (s *Sqlite3Storage) SetExecutedOrder(po *pb.PartialOrder, orderID uint64) error {
	_, err := s.db.Exec(`
		INSERT INTO executed_order(price, amount, book_id, order_id)
		VALUES (?, ?, ?, ?)`,
		po.Price, po.Qty, po.OrderBookId, orderID,
	)
	return err
}

// GetExecutedResults gets the executed order results.
func (s *Sqlite3Storage) GetExecutedResults(orderID uint64) ([]*pb.PartialOrder, error) {
	rows, err := s.db.Query(`
		SELECT price, amount, book_id
		FROM executed_order
		WHERE order_id = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint

	results := make([]*pb.PartialOrder, 0)

	for rows.Next() {
		var po pb.PartialOrder
		err = rows.Scan(&po.Price, &po.Qty, &po.OrderBookId)
		if err != nil {
			return nil, err
		}
		results = append(results, &po)
	}
	return results, nil
}
