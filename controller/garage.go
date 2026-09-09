package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"evertrust-backend-exercise/entity"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type GarageController struct {
	db *sqlx.DB
}

func NewGarageController(db *sqlx.DB) *GarageController {
	return &GarageController{db: db}
}

func (controller *GarageController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type registerReq struct {
		Name            string   `json:"name"`
		Cars            []string `json:"cars"`
		OccupationLimit int      `json:"occupationLimit"`
	}

	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	err := controller.db.GetContext(ctx, &garage, `SELECT id, name, occupation_limit FROM garage WHERE name = ?`, req.Name)
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

	_, err = controller.db.ExecContext(ctx, `INSERT INTO garage(name, occupation_limit) VALUES (?, ?)`, req.Name, req.OccupationLimit)
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
}

func (controller *GarageController) GetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT ME!
}

func (controller *GarageController) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT ME!
}

func (controller *GarageController) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT ME!
}

func (controller *GarageController) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT ME!
}
