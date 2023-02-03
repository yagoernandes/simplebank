package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := 10
	results := make(chan TransferResult, n)
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        int64(amount),
			})
			results <- result
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		result := <-results
		err := <-errs

		require.NoError(t, err)
		require.NotEmpty(t, result)

		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, int64(amount), result.Transfer.Amount)

		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		toEntry := result.ToEntry
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, int64(amount), toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		fromEntry := result.FromEntry
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -int64(amount), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
	}
}
