FROM golang:1.23 AS builder
WORKDIR /app

# Install system dependencies for libheif
RUN apt-get update && apt-get install -y \
    libheif-dev \
    libde265-dev \
    libx265-dev \
    pkg-config \
    build-essential

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o server .

# Runtime - full Debian (NOT distroless)
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libheif1 \
    libde265-0 \
    libx265-199

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
