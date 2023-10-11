.PHONY: build run

build:
	@go build -o ./bin/wserver main.go

run-client: build
	@./bin/wserver -usage=client

run-server: build
	@./bin/wserver -usage=server