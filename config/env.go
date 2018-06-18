package config

import (
	"github.com/caarlos0/env"
)

func load(c interface{}) error {
	return env.Parse(c)
}
