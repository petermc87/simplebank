package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Create a test transfertx func, passing in the usual test args.
func TestTransferTx(t *testing.T) {

	// Create a new db store.
	store := NewStore(testDB)

	// Create two new accounts.
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// Print out balances before.
	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	// Run some concurrent transactions using a loop. It is important to test that concurrency
	// is ok!
	n := 5
	// Specify an amount.
	amount := int64(10)
	// 2. Recieve the errors as one channel.
	errs := make(chan error)

	// 3. Receive the transfer as the second channel.
	results := make(chan TransferTxResult)

	for i := 0; i < 5; i++ {

		// use a go routine. NOTE: There is a set of round brackets at the end to call the func.
		go func() {
			// We are storing this in the context. The context.Background will be passed in as a
			// parent variable.
			// WithValue(parent_Context, key interface{}, val interface{}) --> should not be of string
			// or built in type to avoid collisions.

			ctx := context.Background()
			// Results and err variables that will call the TransferTx function from the store.go file.
			// Go to line 86 to see the breakdown of the function.
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			// 1. Because this is a local go func, we dont have access to require to check for errors
			// So any errors are returned to the main go return. We can use channels

			// 4. Pass data on the right to channel on the left.
			errs <- err
			results <- result

		}()
	}

	// Create a new variable called exisited to check k in the amount and difference checker at the end
	// of this loop.
	// Integer is the key, boolean is the value.
	existed := make(map[int]bool)

	// 5. Check the errors by looping over every new transfer.
	for i := 0; i < n; i++ {
		// 6. Store the errs channel in a new err variable an check for no errors.
		err := <-errs
		require.NoError(t, err)

		// 7. Store result and check if not empty.
		result := <-results
		require.NotEmpty(t, result)

		// 8. Check the transfer object in the result.
		transfer := result.Transfer
		// Not empty
		require.NotEmpty(t, transfer)
		// equal ids and values
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		// not zeros
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Check FromEntry object.
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)

		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// HTTP request. Remember what the query is asking for to complete it.
		// It needs to find the entry by ID
		// NOTE: _ part of the funciton is because we arent storing it in a variable.
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Check ToEntry object.

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)

		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)

		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check from account.
		// Store the result from the FromAccount object.
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// Check to account.
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Print transaction in progress.
		fmt.Println(">> After:", fromAccount.Balance, toAccount.Balance)
		// Check accounts balance.
		// The difference between the balance before
		// account1/account2 and after toAccount and fromAccount.
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		// Both diffs should be equal.
		require.Equal(t, diff1, diff2)
		// If there is a transactions, it will be greater than 0!
		require.True(t, diff1 > 0)
		// If there is a transaction, then the remainder from dividing the amount and the difference should be 0!
		require.True(t, diff1%amount == 0)

		// We are going to have a total of 5 transactions (based in the for loop). So the divider will iterate up:
		// 10/10 = 1, 20/10 = 2, 30/10 = 3 ...

		k := int(diff1 / amount)

		// The number of transactions will determine if this is true or not.
		require.True(t, k >= 1 && k <= n)

		// The map generated should not contain a k value.
		require.NotContains(t, existed, k)

		existed[k] = true

	}

	// Now after balance updates above, check the account balances.
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// Printing balance after transactions.
	fmt.Println("After:", updatedAccount1.Balance, updatedAccount2.Balance)
	// Testing a reduced balance of account1 against the iterated account 1 balance.
	require.Equal(t, account1.Balance-(amount*int64(n)), updatedAccount1.Balance)

	// Testing an increased balance.
	require.Equal(t, account2.Balance+(amount*int64(n)), updatedAccount2.Balance)

}
