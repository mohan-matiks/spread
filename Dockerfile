FROM golang:1.22 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy backend dependencies first
COPY go.mod go.sum ./

RUN go mod download

# Copy all source code
COPY . .

# Build backend server
RUN go build -o server .

# Build frontend assets
WORKDIR /app/web
RUN apt-get update && apt-get install -y nodejs npm
RUN npm install
RUN npm run build

# Create directory for the build files
RUN mkdir -p /app/web/build

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /root/

# Copy backend binary
COPY --from=builder /app/server .

# Copy frontend build files
COPY --from=builder /app/web/build ./web/build

# Set environment variables
ENV SERVE_STATIC=true
ENV STATIC_DIR=/root/web/build

CMD ["./server", "serve"]
