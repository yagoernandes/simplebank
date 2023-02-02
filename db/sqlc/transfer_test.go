package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yagoernandes/simplebank/util"
)

func createRandomTransfer(t *testing.T, account1 Account, account2 Account) Transfer {
	transfer := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	result, err := testQueries.CreateTransfer(context.Background(), transfer)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, transfer.FromAccountID, result.FromAccountID)
	require.Equal(t, transfer.ToAccountID, result.ToAccountID)
	require.Equal(t, transfer.Amount, result.Amount)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	return result
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, account1.ID, transfer2.FromAccountID)
	require.Equal(t, account2.ID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}
