package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

type key int

const dbKey key = 0

// WithDB adds DB to context.
func WithDB(db *sql.DB) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, dbKey, db)
		next(w, r.WithContext(ctx))
	}
}

// DBFromContext gets DB from context.
func DBFromContext(ctx context.Context) *sql.DB {
	return ctx.Value(dbKey).(*sql.DB)
}
