package middleware

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"evertrust-backend-exercise/entity"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type Authenticator struct {
	db *sqlx.DB
}

func NewAuthenticator(db *sqlx.DB) *Authenticator {
	return &Authenticator{db: db}
}

func (middleware *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token-authentication")

		hash := sha256.Sum256([]byte(token))
		var tokenEntity entity.Token

		err := middleware.db.GetContext(r.Context(), &tokenEntity, "SELECT hash, expires_at, created_at FROM token WHERE hash = ?", hash[:])
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("token not found")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				// TODO: RFC 7807 compliant error response
				return
			}
			log.Println(err)
		}
		if tokenEntity.ExpiresAt.Before(time.Now()) {
			log.Println("token expired")
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusUnauthorized)
			// TODO: RFC 7807 compliant error response
			return
		}

		next.ServeHTTP(w, r)
	})
}
