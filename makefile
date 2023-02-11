postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d -p 5432:5432 postgres:15.1-alpine

startdb:
	docker start postgres

createdb:
	docker exec -it postgres createdb -U root -O root simplebank

dropdb:
	docker exec -it postgres dropdb -U root simplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable" -verbose down

sql:
	docker exec -it postgres psql -U root -d simplebank 

sqlc:
	sqlc generate

test:
	go test -v ./... -cover

.PHONY: postgres startdb createdb dropdb migrateup migratedown sql sqlc test