postgres:
	docker run --name postgres-15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine 
createdb:
	docker exec -it postgres-15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-15 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

# test:
# 	go test 
	
.PHONY: postgres createdb dropdb migrateup migratedown sql test