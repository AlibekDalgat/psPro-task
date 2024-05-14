build:
	docker-compose build app

run:
	docker-compose up app

test:
	go test ./internal/delivery/

migrate-up:
	migrate -path ./db -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up
migrate-down:
	migrate -path ./db -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down