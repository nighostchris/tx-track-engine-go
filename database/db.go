package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func Connect(dbUrl string) (db *sql.DB, err error) {
	db, paramsError := sql.Open("postgres", dbUrl)

	if paramsError != nil {
		return nil, fmt.Errorf("[Database Connection] %s", paramsError.Error())
	}

	if error := db.Ping(); error != nil {
		return nil, fmt.Errorf("[Database Connection] Failed to establish database connection: %s", error.Error())
	}

	return db, nil
}

func Migrate(db *sql.DB) string {
	driver, initDriverError := postgres.WithInstance(db, &postgres.Config{})

	if initDriverError != nil {
		return fmt.Sprintf("[Database Migration] Failed to initialize driver: %s", initDriverError.Error())
	} else {
		m, migrationError := migrate.NewWithDatabaseInstance("file://migrations", "blockchain", driver)

		if migrationError != nil {
			return fmt.Sprintf("[Database Migration] Failed to migrate: %s", migrationError.Error())
		} else {
			m.Up()
			return "[Database Migration] Done"
		}
	}
}
