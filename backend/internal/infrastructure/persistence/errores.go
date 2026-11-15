package persistence

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// esViolacionUnicidad detecta el código Postgres 23505 (unique_violation).
func esViolacionUnicidad(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// esViolacionFK detecta el código Postgres 23503 (foreign_key_violation).
func esViolacionFK(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
