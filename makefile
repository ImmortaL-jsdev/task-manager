.PHONY: run build clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/task-manager cmd/server/main.go

clean:
	rm -rf bin/