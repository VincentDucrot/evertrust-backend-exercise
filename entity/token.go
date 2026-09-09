package entity

import "time"

type Token struct {
	Hash      []byte    `db:"hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
