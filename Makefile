include .env
export

run:
	go run ./cmd/server

build:
	go build -o app ./cmd/server

test:
	go test ./...

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)