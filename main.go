package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"evertrust-backend-exercise/config"
	"evertrust-backend-exercise/controller"
	"evertrust-backend-exercise/entity"
	middleware2 "evertrust-backend-exercise/middleware"
	"fmt"
	"log"
	"net/http"

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

	authenticator := middleware2.NewAuthenticator(db)

	tokenController := controller.NewTokenController(db)
	carController := controller.NewCarController(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/tokens", func(r chi.Router) {
		r.Post("/", tokenController.CreateToken)
		r.Get("/{token}", tokenController.GetToken)
		r.Delete("/{token}", tokenController.DeleteToken)
	})
	r.Route("/cars", func(r chi.Router) {
		r.Use(authenticator.Authenticate)

		r.Post("/", carController.Create)
		r.Get("/", carController.GetAll)
		r.Get("/{numberplate}", carController.Get)
		r.Put("/{numberplate}", carController.Update)
		r.Delete("/{numberplate}", carController.Delete)
	})
	r.Route("/garages", func(r chi.Router) {
		r.Use(authenticator.Authenticate)

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			type registerReq struct {
				Name            string   `json:"name"`
				Cars            []string `json:"cars"`
				OccupationLimit int      `json:"occupationLimit"`
			}

			var req registerReq
			if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}

			if req.OccupationLimit <= 0 {
				log.Println("invalid occupationLimit")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}
			if len(req.Cars) > req.OccupationLimit {
				log.Println("cars exceed occupationLimit")
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusBadRequest)
				// TODO: RFC 7807 compliant error response
				return
			}

			var garage entity.Garage
			err = db.GetContext(ctx, &garage, `SELECT id, name, occupation_limit FROM garage WHERE name = ?`, req.Name)
			if err == nil {
				log.Println("garage already registered")
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

			_, err = db.ExecContext(ctx, `INSERT INTO garage(name, occupation_limit) VALUES (?, ?)`, req.Name, req.OccupationLimit)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: RFC 7807 compliant error response
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(req)
		})
	})

	http.ListenAndServe(":3333", r) // TODO: change this to a config
}
