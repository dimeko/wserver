.PHONY: build run

build:
	@go build -o ./bin/ws_server main.go

run: build
	@./bin/ws_server