package controller

import (
	"crypto/sha256"
	"encoding/json"
	"evertrust-backend-exercise/interfaces"
	"evertrust-backend-exercise/service"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type TokenController struct {
	tokenRepo     interfaces.TokenRepository
	tokenDuration time.Duration
}

func NewTokenController(tokenRepo interfaces.TokenRepository, tokenDuration time.Duration) *TokenController {
	return &TokenController{
		tokenRepo:     tokenRepo,
		tokenDuration: tokenDuration,
	}
}

func (controller *TokenController) CreateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plaintext, hash := service.NewToken()
	expiration := time.Now().Add(controller.tokenDuration)

	err := controller.tokenRepo.Create(ctx, hash, expiration)
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
			"expiration": expiration.Format(time.RFC3339),
		},
	)
}

func (controller *TokenController) GetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plaintext := chi.URLParam(r, "token")
	hash := sha256.Sum256([]byte(plaintext))

	token, err := controller.tokenRepo.Get(ctx, hash)
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

	err := controller.tokenRepo.Delete(ctx, hash)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: RFC 7807 compliant error response
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
