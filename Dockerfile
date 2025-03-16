FROM golang:1.22 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server .

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /root/

COPY --from=builder /app/server .

CMD ["./server", "serve"]
