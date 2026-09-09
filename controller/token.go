package controller

import (
	"crypto/sha256"
	"encoding/json"
	"evertrust-backend-exercise/entity"
	"evertrust-backend-exercise/service"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type TokenController struct {
	db *sqlx.DB
}

func NewTokenController(db *sqlx.DB) *TokenController {
	return &TokenController{db: db}
}

func (controller *TokenController) CreateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plaintext, hash := service.NewToken()

	token := entity.Token{
		Hash:      hash[:],
		ExpiresAt: time.Now().Add(time.Minute * 10), // TODO: change this to a config
	}

	_, err := controller.db.ExecContext(ctx, `INSERT INTO token(hash, expires_at) VALUES (?, ?)`, token.Hash, token.ExpiresAt)
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
}

func (controller *TokenController) GetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plaintext := chi.URLParam(r, "token")
	hash := sha256.Sum256([]byte(plaintext))
	token := entity.Token{}

	err := controller.db.GetContext(ctx, &token, "SELECT hash, expires_at, created_at FROM token WHERE hash = ?", hash[:])
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
}

func (controller *TokenController) DeleteToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plaintext := chi.URLParam(r, "token")
	hash := sha256.Sum256([]byte(plaintext))

	_, err := controller.db.ExecContext(ctx, `DELETE FROM token WHERE hash=?`, hash[:])
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: RFC 7807 compliant error response
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
