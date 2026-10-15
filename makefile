.PHONY: run build clean lint docker-up docker-down test test-unit test-integration test-cover

run:
	go run cmd/server/main.go

build:
	go build -o bin/task-manager cmd/server/main.go

clean:
	rm -rf bin/

lint:
	golangci-lint run ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

test:
	go test ./... -v

test-unit:
	go test ./internal/service/... -v

test-integration:
	TESTCONTAINERS_RYUK_DISABLED=true go test ./internal/repository/... -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out