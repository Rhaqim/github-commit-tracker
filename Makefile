.PHONY: setup build run test

setup:
	cp .env.example .env && go mod tidy

build:
	go build -o main cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./...
