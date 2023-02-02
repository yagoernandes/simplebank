// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"
)

type Account struct {
	ID        int64
	Owner     string
	Currency  string
	Balance   int64
	CreatedAt time.Time
}

type Entry struct {
	ID        int64
	AccountID int64
	// can be positive or negative
	Amount    int64
	CreatedAt time.Time
}

type Transfer struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	// must be positive
	Amount    int64
	CreatedAt time.Time
}
