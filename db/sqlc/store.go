package db

import (
	"context"
	"database/sql"
	"fmt"
)


// To execute all functions and transactions.
type Store struct {
	// So we can create a combination of queries.
	*Queries
	// To create a new DB transaction, we need the DB object.
	db * sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}


// Create pass in context and a call back function. Once a queries object is created, its passed into
// the callback function.
// execTx is a function execution within a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Returns a transaction OR err.
	tx, err := store.db.BeginTx(ctx, nil)

	// If the err exists, return the object.
	if err != nil {
		return err
	}

	// Creating a new transaction query.
	q := New(tx)

	// Callback function using Queries (see line 23 where it says fn func(...))
	err = fn(q)

	// If there is an error, execure the rollback func.
	if err != nil {
		// If the rollback has a error, two errors are reported.
		if rbErr := tx.Rollback(); rbErr != nil {
			// The errored message will take in the err variable and rberr variable and print them.
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	// If successful, commit the transaction.
	return tx.Commit()
}

// Transfer transaction parameters
// from, to, results
// var, in64 json

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id`
	Amount int64 `json:"amount`
}



// Transfer transaction results
// from to Account json
// From to entry json
// Trnasfer transfer json

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json: "from_account"`
	ToAccount Account `json: "to_account`
	FromEntry Entry `json: "from_entry`
	ToEntry Entry `json: "to_entry`
}

// TransferTx is a money transfer database entry from one account to another.
// This is achieved via a transfer record, account entries and updates to transaction balance.
func (store *Store) TransferTx (ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// NOTE: Ctx means everything else besides the arguements passed onto the function. This helps with
	// understandin errors if a connection goes down, for example.



	// Passing in the queries object as an arg to the callback.
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Accessing the 'result' from the transfertxresults function on line 79 using the 
		// func(q ...) callback on line 86. This will create a closure function. Similarly with arg.
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams {
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		// To define the from entry, we need to store the account ID as the from account from the tx.
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams {
			AccountID: arg.FromAccountID,
			// Negative amount because money is being deducted.
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams {
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: Update account balance.
		// Get the account from the database.
		// Update amount.
		// Send it back to database.

		
		// Getting account based on the from account id.
		account1, err := q.GetAccount(ctx, arg.FromAccountID)

		if err != nil {
			return err
		}

		// Updating the balance based on the arguements passed in.
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.FromAccountID, 
			// Negating the transfered amount.
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		account2, err := q.GetAccount(ctx, arg.ToAccountID)

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}