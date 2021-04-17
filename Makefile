chaton:
	@echo Chaton: Simple chat service using gRPC.

################################################## Server
# Run server localy
server-run:
	@go run ./server/*

# Build docker image
server-image:
	@docker build --no-cache -t chaton-server -f ./images/server.Dockerfile .

# Run container
server-ctn-run:
	@if [ -f ./server/.env ]; then \
		docker run --env-file ./server/.env -d --name chaton-server -p 21617:21617 chaton-server; \
	else \
		docker run -d --name chaton-server -p 21617:21617 chaton-server; \
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
# Run client localy
client-run:
	@go run ./clients/*

# Build the client
client-build:
	@mkdir -p ./bin
	@go build -o ./bin/chaton ./clients/*
	@echo "Client build (./bin/chaton)"
