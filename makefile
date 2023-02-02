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

sqlc:
	sqlc generate

.PHONY: postgres startdb createdb dropdb migrateup migratedown sqlc