package main

import (
	"crypto/sha256"
	"encoding/json"
	"evertrust-backend-exercise/config"
	"evertrust-backend-exercise/entity"
	"evertrust-backend-exercise/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s?parseTime=true", dbConfig.Username, dbConfig.Password, dbConfig.Url))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/tokens", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			plaintext, hash := service.NewToken()

			token := entity.Token{
				Hash:      hash[:],
				ExpiresAt: time.Now().Add(time.Minute * 10), // TODO: change this to a config
			}

			_, err := db.ExecContext(ctx, `INSERT INTO token(hash, expires_at) VALUES (?, ?)`, token.Hash, token.ExpiresAt)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(
				map[string]string{
					"token":      plaintext,
					"expiration": time.Now().Add(time.Minute * 10).Format(time.RFC3339),
				},
			)
		})
		r.Get("/{token}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			plaintext := chi.URLParam(r, "token")
			hash := sha256.Sum256([]byte(plaintext))
			token := entity.Token{}

			err := db.GetContext(ctx, &token, "SELECT hash, expires_at, created_at FROM token WHERE hash = ?", hash[:])
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}
			if token.ExpiresAt.Before(time.Now()) {
				log.Println("Token expired")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusUnauthorized)
				// TODO: RFC 7807 compliant error response
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(
				map[string]string{
					"token":      plaintext,
					"expiration": token.ExpiresAt.Format(time.RFC3339),
				},
			)
		})
	})

	http.ListenAndServe(":3333", r) // TODO: change this to a config
}
