.PHONY: chaton

SERVER_PORT ?= 21617

chaton:
	@echo Chaton: Simple chat service using gRPC.

################################################## Server
.PHONY: server-run server-image server-ctn-run server-ctn-restart server-ctn-stop server-ctn-delete

# Run server localy
server-run:
	@go run ./server/*

# Build docker image
server-image:
	@docker build --no-cache -t chaton-server -f ./images/server.Dockerfile .

# Run container
server-ctn-run:
	@if [ -f ./server/.env ]; then \
		docker run --env-file ./server/.env -d --name chaton-server -p $(SERVER_PORT):$(SERVER_PORT) chaton-server; \
	else \
		docker run -d --name chaton-server -p $(SERVER_PORT):$(SERVER_PORT) chaton-server; \
	fi

# Restart the running container with a new builded version
server-ctn-restart: server-ctn-stop server-ctn-delete server-image server-ctn-run
	@echo Container stopped, deleted, image rebuilded, container now running.

# Stop contaner
server-ctn-stop:
	@docker stop chaton-server

# Delete container
server-ctn-delete:
	@docker container rm chaton-server

################################################## Go cli client
.PHONY: go-cli-client-run go-cli-client-build

# Run client localy
go-cli-client-run:
	@go run ./clients/go_cli_client/*

# Build the client
go-cli-client-build:
	@mkdir -p ./bin
	@go build -o ./bin/chaton ./clients/go_cli_client/*
	@echo "Client build (./bin/chaton)"
