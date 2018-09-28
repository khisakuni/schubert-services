package config

// DB is a struct that holds configs for database.
type DB struct {
	User     string `env:"DB_USER" envDefault:"koheihisakuni"`
	Password string `env:"DB_PASSWORD" envDefault:"shiba"`
	Name     string `env:"DB_NAME" envDefault:"schubert"`
	SSLMode  string `env:"DB_SSL" envDefault:"disable"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
}

// Load loads the configs from the environment.
func (c *DB) Load() error {
	return load(c)
}
