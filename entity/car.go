package entity

import "time"

type Car struct {
	Numberplate string    `db:"numberplate"`
	Model       string    `db:"model"`
	Color       string    `db:"color"`
	Serial      string    `db:"serial"`
	CreatedAt   time.Time `db:"created_at"`
}
