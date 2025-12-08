FROM golang:1.22

WORKDIR /app

# Install libheif + deps for HEIC conversion
RUN apt-get update && apt-get install -y \
    libde265-dev \
    libheif-dev && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o server .

CMD ["./server"]
