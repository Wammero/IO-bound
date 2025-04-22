FROM golang:1.23.1-bullseye

WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1

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
    wget \
 && wget https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
 && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
 && rm dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
 && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 go build -o server ./cmd/app

CMD ["dockerize", "-wait", "tcp://kafka:9092", "-wait", "tcp://postgres:5432", "-timeout", "60s", "./server"]
