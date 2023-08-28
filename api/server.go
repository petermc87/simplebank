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
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}
