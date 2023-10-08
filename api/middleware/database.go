package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

// Middleware function to inject the db variable into the request context
func WithDatabase(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a context with the db variable
		ctx := context.WithValue(r.Context(), "db", db)

		// Call the next handler with the updated request context
		next(w, r.WithContext(ctx))
	}
}
