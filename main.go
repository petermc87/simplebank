package main

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// func main() {
// 	conn, err := sql.Open(dbDriver, dbSource)
// 	if err != nil {
// 		log.fatal("cannot connect to db", err)
// 	}
// 	store := db.NewStore(conn)
// 	server := api.NewServer(store)
// }
