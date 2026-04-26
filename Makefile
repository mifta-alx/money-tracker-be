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

migrate-down-all:
	migrate -path migrations -database "$(DB_URL)" down

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

dev:
	docker-compose --env-file .env -f docker/docker-compose.yaml up --build -d

stop:
	docker-compose -f docker/docker-compose.yaml down

logs:
	docker logs -f money-tracker-api