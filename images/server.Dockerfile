FROM golang:1.18 as builder

WORKDIR /usr/src/chaton

COPY ./server ./server
COPY ./grpc ./grpc
COPY go.mod .
COPY go.sum .

WORKDIR /usr/src/chaton/server

RUN CGO_ENABLED=0 go build -o /bin/server


FROM alpine

WORKDIR /bin/chaton

COPY --from=builder /bin/server .

CMD ["./server"]
