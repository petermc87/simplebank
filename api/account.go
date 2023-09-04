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
