package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
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
		r.Delete("/{token}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			plaintext := chi.URLParam(r, "token")
			hash := sha256.Sum256([]byte(plaintext))

			_, err := db.ExecContext(ctx, `DELETE FROM token WHERE hash=?`, hash[:])
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})
	})
	r.Route("/cars", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			type registerReq struct {
				Numberplate string `json:"numberplate"`
				Model       string `json:"model"`
				Color       string `json:"color"`
				Serial      string `json:"serial"`
			}

			var req registerReq

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}

			var car entity.Car

			err := db.GetContext(ctx, &car, `SELECT numberplate, model, color, serial FROM car WHERE numberplate = ?`, req.Numberplate)
			if err == nil {
				log.Println("car already registered")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusConflict)
				// TODO: RFC 7807 compliant error response
				return
			}
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}

			_, err = db.ExecContext(ctx, `INSERT INTO car(numberplate, model, color, serial) VALUES (?, ?, ?, ?)`, req.Numberplate, req.Model, req.Color, req.Serial) // TODO: Make serial unique?
			if err != nil {
				log.Println("car already registered")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(req)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			type listCarsResp struct {
				Numberplate string `json:"numberplate"`
				Model       string `json:"model"`
				Color       string `json:"color"`
				Serial      string `json:"serial"`
			}

			var cars entity.Cars

			err = db.SelectContext(ctx, &cars, `SELECT numberplate, model, color, serial FROM car`)
			if err != nil {
				log.Println("car already registered")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			// TODO: Make function
			resp := make([]listCarsResp, len(cars))
			for i, car := range cars {
				resp[i] = listCarsResp{
					Numberplate: car.Numberplate,
					Model:       car.Model,
					Color:       car.Color,
					Serial:      car.Serial,
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})
		r.Get("/{numberplate}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			numberplate := chi.URLParam(r, "numberplate")
			var car entity.Car

			type carDetailResp struct {
				Numberplate string `json:"numberplate"`
				Model       string `json:"model"`
				Color       string `json:"color"`
				Serial      string `json:"serial"`
			}

			err = db.GetContext(ctx, &car, `SELECT numberplate, model, color, serial FROM car WHERE numberplate = ?`, numberplate)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Println("car not found")
					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(http.StatusBadRequest)
					// TODO: RFC 7807 compliant error response
					return
				}
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(carDetailResp{
				Numberplate: car.Numberplate,
				Model:       car.Model,
				Color:       car.Color,
				Serial:      car.Serial,
			})
		})
		r.Put("/{numberplate}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			numberplate := chi.URLParam(r, "numberplate")
			if numberplate == "" {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}

			type updateCarReq struct {
				Model  string `json:"model"`
				Color  string `json:"color"`
				Serial string `json:"serial"`
			}
			var req updateCarReq

			if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}

			tx, err := db.BeginTxx(ctx, nil)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}
			defer tx.Rollback()

			var car entity.Car

			err = tx.GetContext(ctx, &car, `SELECT numberplate, model, color, serial FROM car WHERE numberplate = ?`, numberplate)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Println("car not found")
					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(http.StatusBadRequest)
					// TODO: RFC 7807 compliant error response
					return
				}
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			_, err = tx.ExecContext(ctx, `UPDATE car SET model=?, color=?, serial=? WHERE numberplate=?`, req.Model, req.Color, req.Serial, numberplate)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			tx.Commit()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"numberplate": car.Numberplate,
				"model":       req.Model,
				"color":       req.Color,
				"serial":      req.Serial,
			})
		})
	})

	http.ListenAndServe(":3333", r) // TODO: change this to a config
}
