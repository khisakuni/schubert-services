package database

import (
	"database/sql"
	"fmt"

	"github.com/khisakuni/schubert-services/config"

	// Load db driver
	_ "github.com/lib/pq"
)

func New() (*sql.DB, error) {
	c := config.DB{}
	c.Load()
	dbstring := fmt.Sprintf("user=%s dbname=%s sslmode=%s", c.User, c.Name, c.SSLMode)
	return sql.Open("postgres", dbstring)
}
