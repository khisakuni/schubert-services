package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE users(id SERIAL PRIMARY KEY NOT NULL, email VARCHAR (255), password VARCHAR (128), username VARCHAR (20));")
	return err
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE users;")
	return err
}
