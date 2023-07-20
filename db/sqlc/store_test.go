package db

import (
	"testing"
	"context"

	"github.com/stretchr/testify/require"
)

// Create a test transfertx func, passing in the usual test args.
func TestTransferTx(t *testing.T){

	// Create a new db store.
	store := NewStore(testDB)

	// Create two new accounts.
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Specify an amount.
	amount := int64(10)

	// Run some concurrent transactions using a loop. It is important to test that concurrency
	// is ok!
	n := 5

	// 2. Recieve the errors as one channel.
	errs := make(chan error)

	// 3. Receive the transfer as the second channel.
	results := make(chan TransferTxResult)

	for i := 0; i < 5; i ++ {
		// use a go routine. NOTE: There is a set of round brackets at the end to call the func.
		go func(){
			// Results and err variables that will call the TransferTx function from the store.go file.
			// Go to line 86 to see the breakdown of the function.	
			result, err := store.TransferTx(context.Background(), TransferTxParams {
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			// 1. Because this is a local go func, we dont have access to require to check for errors
			// So any errors are returned to the main go return. We can use channels 


			// 4. Pass data on the right to channel on the left.
			errs <- err
			results <- result

		}()
	}
		
	// 5. Check the errors by looping over every new transfer.
	for i := 0; i < n; i ++ {
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




		// // Check FromEntry object.
		// fromEntry := result.FromEntry

		// require.NotEmpty(t, fromEntry)

		// require.Equal(t, fromEntry.AccountID, account1.ID)
		// require.Equal(t, fromEntry.amount, -amount)
		
		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)

		// // HTTP request. Remember what the query is asking for to complete it.
		// // It needs to find the entry by ID
		// // NOTE: _ part of the funciton is because we arent storing it in a variable.
		// _, err := store.GetEntry(context.Background(), fromEntry.ID)
		// require.NoError(t, err)


		// Check ToEntry object.

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)

		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)

		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts balance.
	
	}


}
