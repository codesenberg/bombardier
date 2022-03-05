FROM golang:1.15 as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY *.go/ ./
COPY internal/ internal/

ARG VERSION="0.0.1"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
    go build \
    -ldflags "-X main.Version=$VERSION" \
    -a \
    -o bombardier-linux-amd64

FROM debian:bookworm-slim

WORKDIR /

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN addgroup --gid 901 bombardier && \
    adduser --uid 901 --gid 901 bombardier

USER bombardier

COPY --from=builder /build/bombardier-linux-amd64 .

ENTRYPOINT ["./bombardier-linux-amd64"]
