.PHONY: run build clean lint

run:
	go run cmd/server/main.go

build:
	go build -o bin/task-manager cmd/server/main.go

clean:
	rm -rf bin/

lint:
	golangci-lint run ./...