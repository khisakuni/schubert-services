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
	dbstring := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s", c.Host, c.User, c.Name, c.SSLMode, c.Password)
	return sql.Open("postgres", dbstring)
}
