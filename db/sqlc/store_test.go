package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	var amount int64 = 10
	results := make(chan TransferResult, n)
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
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

	// Check the new account balances
	totalAmount := int64(n) * amount

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	require.Equal(t, account1.Balance-totalAmount, updatedAccount1.Balance)

	var updatedAccount2 Account
	updatedAccount2, err = store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance+totalAmount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	var amount int64 = 10
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			fromAccountID := account1.ID
			toAccountID := account2.ID

			if i%2 == 1 {
				fromAccountID = account2.ID
				toAccountID = account1.ID
			}
			fmt.Printf("Transfer %d from %d to %d\n", i, fromAccountID, toAccountID)
			_, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}(i)
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the new account balances

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	var updatedAccount2 Account
	updatedAccount2, err = store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
