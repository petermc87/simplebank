package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func main() {
	// Taking in the env variables from the app.env.
	config, err := util.LoadConfig(".")

	// Handling a load config error.
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	// Creating a sql connection using the config settings.
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	// Create a new connection and store.
	store := db.NewStore(conn)
	server := api.NewServer(store)

	// Error handling.
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
