package entity

import "time"

type Garage struct {
	Id              int       `db:"id"`
	Name            string    `db:"name"`
	OccupationLimit int       `db:"occupation_limit"`
	CreatedAt       time.Time `db:"created_at"`
}

type GarageCar struct {
	GarageId  int `db:"garage_id"`
	CarId     int `db:"car_id"`
	CreatedAt time.Time
}
