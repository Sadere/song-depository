// Provides necessary functions for communicating with DB
package database

import (
	"github.com/jmoiron/sqlx"
)

// Returns instance of DB connection
func NewConnection(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
