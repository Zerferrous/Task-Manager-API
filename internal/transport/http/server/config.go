package server

import (
	"time"
)

type Config struct {
	Addr            string        `envconfig:"HTTP_ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" required:"true"`
}
