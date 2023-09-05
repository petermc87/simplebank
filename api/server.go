package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

// HTTP service data structure.
type Server struct {
	// This is bases on it on the db connect made in the store.go
	store *db.Store
	// Router and handler standard method in Gin.
	router *gin.Engine
}

// New server creates a new instance of a server object where The
// routing and HTTP verbs are defined.
func NewServer(store *db.Store) *Server {
	// store is the input.
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)

	router.GET("/accounts", server.listAccount)

	router.PUT("/accounts/:id", server.updateAccount)

	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on a specefied address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// It will take in a gin object -> map string interface
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
