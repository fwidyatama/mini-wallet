postgres:
	docker run --name postgres-local -p 5433:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -d postgres

createdb:
	docker exec -it postgres-local createdb --username=postgres --owner=postgres mini_wallets

dropdb:
	docker exec -it postgres-local dropdb mini_wallets --username=postgres

migrateup:
	 migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5433/mini_wallets?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5433/mini_wallets?sslmode=disable" -verbose down

dev:
	go run cmd/main.go


.PHONY: postgres createdb dropdb migrateup migratedown dev