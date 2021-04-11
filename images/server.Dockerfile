FROM golang:1.16 as builder

WORKDIR /usr/src/chaton

COPY ./server ./server
COPY ./grpc ./grpc
COPY go.mod .

RUN go mod download

WORKDIR /usr/src/chaton/server

RUN CGO_ENABLED=0 go build -o /bin/server


FROM alpine

WORKDIR /bin/chaton

COPY --from=builder /bin/server .

CMD ["./server"]
