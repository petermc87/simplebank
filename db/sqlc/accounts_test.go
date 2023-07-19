// The same package as the querying code.
package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"

)

//Create a random account.
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
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

	return account
}

// Call the function create account for this test case.
func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}


func TestGetAccount(t *testing.T) {
	// First create an account
	account1 := createRandomAccount(t)
	// Then second query the account by the first accounts ID and store it as account2
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	// Use the same requires to check if its good to be passed.
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// Check id and othe peramters match for both records.
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	// You can also check if they are within a certain time of each other.
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: 	 account1.ID,
		Balance: util.RandomMoney(),
	}

	// Updating the account based on the inputted arguements.
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	// We compare what the args inputted by the user are to the found account (account2)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	// Passing in the id to be deleted. We dont need to store a secon account
	// because we are not doing a comparison.
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// Checking if the account exists after deletion.
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// Checking if it returns an error.
	require.Error(t, err)
	// This checks that there are no rows related to this entry.
	require.EqualError(t, err, sql.ErrNoRows.Error())
	// Checking that account 2 is an empty variable
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	// Creating multiple accounts.
	for i := 0; i < 10; i ++ {
		createRandomAccount(t)
	}

	// The params is that we skip the first accounts and then check after that.
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	// Checking the length is 5 for the offset portion.
	require.Len(t, accounts, 5)


	// Looping through the list and checking to see if each account is not empty.
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}