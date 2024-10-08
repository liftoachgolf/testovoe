postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
creatdedb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path schema/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path schema/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
	
migrateInit:
	migrate create -ext sql -dir db/migrations -seq init       