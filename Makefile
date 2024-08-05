.PHONY: build run test

build:
	go build -o main .

run:
	go run cmd/main.go

test:
	go test ./...
