package db

import (
	"github.com/jackc/pgx/v4"
)

// PG is a postgres database driver
type PG struct {
	connection *pgx.Conn
}

// NewPG creates a new postgres database driver
func NewPG(c *pgx.Conn) PG {
	return PG{
		connection: c,
	}
}
