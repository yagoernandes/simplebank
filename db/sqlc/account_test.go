package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yagoernandes/simplebank/util"
)

func TestCreateAccount(t *testing.T) {
	acc := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), acc)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, acc.Owner, account.Owner)
	require.Equal(t, acc.Balance, account.Balance)
	require.Equal(t, acc.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
