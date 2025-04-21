FROM golang:1.23.1-bullseye

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git \
    gcc \
    g++ \
    make \
    pkg-config \
    zlib1g-dev \
    libssl-dev \
    libsasl2-dev \
    libzstd-dev \
    libcurl4-openssl-dev \
    liblz4-dev \
    librdkafka-dev \
 && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 go build -o server ./cmd/app

CMD ["./server"]
