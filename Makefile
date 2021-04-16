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
server-run-ctn:
	@docker run -d --name chaton-server chaton-server

# Stop contaner
server-stop-ctn:
	@docker stop chaton-server

################################################## Go cli client
# Run client localy
client-run:
	@go run ./clients/*
