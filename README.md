# Chaton

Simple chat service using gRPC.

## The service

The gRPC service is described in `./proto/chaton.proto`

## The server

It's implemented in Go and can be build and executed with the `Makefile`.

### Configuration

The defaults address and port the server listen is `0.0.0.0:21617`.  
It can be configures with environement variable `HOST` and `PORT`.

You can use the `.env` file (to complete with the `.env.tpl`) for docker usage.

## Clients

There is multiple clients in the `./clients` directory.
