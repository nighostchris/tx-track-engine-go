package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseConnectionParams struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func Connect(params DatabaseConnectionParams) (db *sql.DB, err error) {
	var connectionUrl string

	if params.Port != "" {
		connectionUrl = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			params.Username,
			params.Password,
			params.Host,
			params.Port,
			params.Database,
		)
	} else {
		connectionUrl = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			params.Username,
			params.Password,
			params.Host,
			params.Database,
		)
	}

	db, paramsError := sql.Open("postgres", connectionUrl)

	if paramsError != nil {
		return nil, fmt.Errorf("[Database Connection] %s", paramsError.Error())
	}

	if error := db.Ping(); error != nil {
		return nil, fmt.Errorf("[Database Connection] Failed to establish database connection: %s", error.Error())
	}

	return db, nil
}
