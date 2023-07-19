package db

import (
	"context"
	// "database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"

)


// Create a random entry.
func createRandomEntry(t *testing.T, account Account) Entry {
	// Pass in the account ID and make it a parameter.
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry

}

// We create a new entry with an assocated random account creation.
func TestCreateEntry(t *testing.T) {
	// Account gets passed in as a variable to the create entry function.
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}


func TestGetEntry(t *testing.T) {
	// First get a random account to put an accountID to.
	account := createRandomAccount(t)
	// Create a new entry and pass in account.
	entry := createRandomEntry(t, account)

	// Perform the get operation using the ID.
	entryResponse, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryResponse)

	// Checking if different elements within the object match.
	require.Equal(t, entry.ID, entryResponse.ID)
	require.Equal(t, entry.AccountID, entryResponse.AccountID)
	require.Equal(t, entry.Amount, entryResponse.Amount)
	require.WithinDuration(t, entry.CreatedAt, entryResponse.CreatedAt, time.Second)
}


