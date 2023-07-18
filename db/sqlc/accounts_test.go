// The same package as the querying code.
package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

)

// The *testing.T is the object being passed in. It is always uppercase t for testing.
func TestCreateAccount(t *testing.T){
	arg := CreateAccountParams{
		Owner: "tom",
		Balance: 100,
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	//Importing the testify module to check different entries sent to the table.
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
