package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferResult, error) {
	var transferResult TransferResult
	var err error

	err = store.execTx(ctx, func(q *Queries) error {
		transferResult.Transfer, err = q.CreateTransfer(ctx, arg)
		if err != nil {
			return err
		}

		transferResult.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			fromAccount, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

			transferResult.FromAccount = fromAccount

			toAccount, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

			transferResult.ToAccount = toAccount

			_, err = q.AddBalance(ctx, AddBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			_, err = q.AddBalance(ctx, AddBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			toAccount, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

			transferResult.ToAccount = toAccount

			fromAccount, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

			transferResult.FromAccount = fromAccount

			_, err = q.AddBalance(ctx, AddBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}

			_, err = q.AddBalance(ctx, AddBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	return transferResult, err
}

// analyzing this code the deadlock is possible because of the following:
// 1. the first goroutine locks the account with id 1 and then tries to lock the account with id 2
// 2. the second goroutine locks the account with id 2 and then tries to lock the account with id 1
// 3. the first goroutine tries to lock the account with id 2 and waits for the second goroutine to unlock it
// 4. the second goroutine tries to lock the account with id 1 and waits for the first goroutine to unlock it
// 5. the deadlock happens
