postgres:
	docker run --name postgres_soundproff -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres_soundproff createdb --username=root --owner=root soundproof_db

dropdb:
	docker exec -it postgres_soundproff dropdb soundproof_db

swagger:
	swag init -dir ./cmd,./internal/transport/http/

lint-host: ## Run golangci-lint directly on host
	@echo "> Linting..."
	golangci-lint run -c .golangci.yml -v
	@echo "> Done!"

start:
	make postgres
	make createdb
	make start server

start-server:
	go run cmd/main.go  

migrateup:
	migrate -path ./db/migration -database "postgresql://root:postgres@$(DB_HOST):5432/soundproof_db?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -database "postgresql://root:postgres@$(DB_HOST):5432/soundproof_db?sslmode=disable" -verbose down

mock:
	mockgen -package mock -destination internal/domain/model/mock/mock.go -source=internal/domain/model/domain.go
	mockgen -package mock -destination internal/domain/service/mock/service_mock.go -source=internal/domain/service/service.go

test:
	go test ./... -v -coverpkg=./...

.PHONY: postgres createdb dropdb start migrateup start-server lint-host mock test swagger
