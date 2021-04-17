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
		docker run --env-file ./server/.env -d --name chaton-server chaton-server; \
	else \
		docker run -d --name chaton-server chaton-server; \
	fi


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
