package db

import (
	"database/sql"
	"fmt"
	log "github.com/binpossible49/go-libs/log"

	"go.uber.org/zap"
)

type dbHelper struct {
	db *sql.DB
}

// NewDBHelper creates an instance
func NewOracleDBHelper(host string, port int, username, password, database string) DBHelper {
	db, err := initOracle(host, port, username, password, database)
	if err != nil {
		log.Logger.Panic("Failed to init oracle", zap.Error(err))
	}
	return &dbHelper{
		db: db,
	}
}

func (h *dbHelper) Open() *sql.DB {
	return h.db
}

func (h *dbHelper) Close() error {
	return h.db.Close()
}

func (h *dbHelper) Begin() (*sql.Tx, error) {
	return h.db.Begin()
}

func (h *dbHelper) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (h *dbHelper) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func initOracle(host string, port int, username, password, database string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%v/%v@%v:%v/%v", username, password, host, port, database)

	db, err := sql.Open("oci8", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
