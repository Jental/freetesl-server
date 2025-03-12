package db

import (
	"database/sql"

	"github.com/jental/freetesl-server/common"
	_ "github.com/lib/pq"
)

func OpenAndTestConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", common.DB_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}
