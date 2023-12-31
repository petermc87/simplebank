package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

// The data type for the object being created.
type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	// Go to gin and look up binding for more on how to be more specific. -> binding to JSON data
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

// If you go to the POST request in the server.go file and hover over it,
// it will redirect you to the breakdown of whats required. Here, one of the requirements
// is the context object.
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// The context should be outputted to the screen as JSON.
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Passing in params from the req body.
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	// Creating the account in the database.
	account, err := server.store.CreateAccount(ctx, arg)

	// Sending JSON data and output to the client.
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
	}

	// If there is no error, send status ok JSON to the terminal and account to the client.
	ctx.JSON(http.StatusOK, account)

}

// Get account request takes in the id through the URI gin parameter.
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		// If the ID doesnt exist
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
		}
		// General error
		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
	}

	ctx.JSON(http.StatusOK, account)

}

// Get account list using bind qeury params fron gin: https://gin-gonic.com/docs/examples/only-bind-query-string/
// We want a minimum page size of results.
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We need to declare limit and offset params.ctx
	arg := db.ListAccountsParams{
		// The limit is the page size.
		Limit: req.PageSize,
		// Off set needs to be calculated.
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	// Error responses
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Account returned.
	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	ID      int64 `uri:"id" binding:"required,min=1"`
	Balance int64 `json:"balance"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}
	account, err := server.store.UpdateAccount(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

	}

	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, gin.H{"response": "status Ok"})
}
