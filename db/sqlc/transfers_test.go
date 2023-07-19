package db

import (
	"time"
	"context"
	"testing"

	"github.com/techschool/simplebank/util"
	"github.com/stretchr/testify/require"
)

// Create a random transfer.
func CreateRandomTransfer(t *testing.T, account1, account2 Account) Transfer {

	// Define the arguements for the function
	arg := CreateTransferParams {
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomMoney(),
	}

	// Create and store the function
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	// Test for error, empty, equal, notZero (would be the id and the created at)	
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}



// Create a function that calls the random transfer and passed in two random accounts.
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Call the function.
	CreateRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	// Create two random accounts.
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Create a new transfer.
	transfer := CreateRandomTransfer(t, account1, account2)

	// Perform the HTTP get request.
	transferResponse, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	// Test for Noterror, Not empty, 
	require.NoError(t, err)
	require.NotEmpty(t, transferResponse)
	
	// if each entry equal each other, and within duration for created at.
	require.Equal(t, transfer.FromAccountID, transferResponse.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferResponse.ToAccountID)
	require.Equal(t, transfer.Amount, transferResponse.Amount)
	require.WithinDuration(t, transfer.CreatedAt, transferResponse.CreatedAt, time.Second)

}


// Create ListTransfer.
func TestListTransfer(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t) 
	// Create a loop of 10 that will pass in two created accounts.
	for i := 0; i < 5; i ++ {
		CreateRandomTransfer(t, account1, account2)
		CreateRandomTransfer(t, account2, account1)
	}

	// The args for the list entry params are for the offset and limit, and either from or to account.
	// NOTE: These are the args that are required to complete the SQL query. Check the query file for each $ 
	// value.

	arg := ListTransfersParams {
		// We only need to check onte of the accounts created in either the from or to accounts because
		// we specified an OR in the Postgres query.
		FromAccountID: account1.ID,
		ToAccountID: account1.ID,
		Limit: 5,
		Offset: 5,
	}

	// Create and store the list.
	transferList, err := testQueries.ListTransfers(context.Background(), arg)

	// Check for no errors and the len of the transfers list.
	require.NoError(t, err)
	require.Len(t, transferList, 5)

	// Check each entry in a loop within the range of the entries list for:
	for _, transfer := range transferList{
		// NotEmpty
		require.NotEmpty(t, transfer)
		// Checking if its true the current transfer has a to or from param with 
		// the matching ID of the first account. ONly need to check on account because there
		// are always going to be two accounts attached to a transfer.
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)

	}
}


