package entity

import "time"

type Car struct {
	Id          int       `db:"id"`
	Numberplate string    `db:"numberplate"`
	Model       string    `db:"model"`
	Color       string    `db:"color"`
	Serial      string    `db:"serial"`
	CreatedAt   time.Time `db:"created_at"`
}

type Cars []Car
