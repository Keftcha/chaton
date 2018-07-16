#!/usr/bin/env python3

import argparse
import socket
import select
import time

parser = argparse.ArgumentParser()

# Declare the arguments required
parser.add_argument(
        "--host",
        type=str,
        help="Please, give a hostname or a IP address"
    )
parser.add_argument(
        "--port",
        type=int,
        help="Please, give a port number"
    )
parser.add_argument(
        "--clients",
        type=int,
        help="Please, give the maximum clients connected"
    )
arguments = parser.parse_args()

# Check if all arguments are given
host = arguments.host if arguments.host else "localhost"
port = arguments.port if arguments.port else 21617
max_clients = arguments.clients if arguments.clients else 10

# Define server parameters
server = socket.socket()
server.bind((host, port))
server.listen(max_clients)

print("The server is listening on port {0}\n".format(port))

running = True
clients = []    # List of all clients

while running:
    # Accept connection of socket
    connecting_requests, _, _ = select.select([server], [], [], 0)
    for request in connecting_requests:
        client, info_client = request.accept()
        clients.append(client)

    if clients:
        rlist, _, _ = select.select(clients, [], [], 0)
        for client_talking in rlist:
            message = client_talking.recv(4096)
            text = message.decode()[:-1]    # Remove the "\n" at the end

            # Send message to all other connected clients
            [client.send(message) for client in clients]

            print(time.strftime("%Y/%m/%d %H:%M:%S"))
            print(str(len(clients)) + " connected sockets")
            print("message send: " + text + "\n")

            if text == "leave()":
                client_talking.close()
                del clients[clients.index(client_talking)]
