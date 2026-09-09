package main

import (
	"evertrust-backend-exercise/config"
	"evertrust-backend-exercise/controller"
	"evertrust-backend-exercise/datastore/mysql"
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
	generalConfig := config.LoadGeneralConfig()

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s?parseTime=true", dbConfig.Username, dbConfig.Password, dbConfig.Url))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authenticator := middleware2.NewAuthenticator(db)

	tokenRepo := mysql.NewTokenRepo(db)

	tokenController := controller.NewTokenController(tokenRepo, generalConfig.TokenDuration)
	carController := controller.NewCarController(db)
	garageController := controller.NewGarageController(db)

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

		r.Post("/", garageController.Create)
	})

	log.Printf("Server is running on port %s\n", generalConfig.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", generalConfig.Port), r)
}
