FROM golang:1.18-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o grpc-echo .

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.7 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /app/grpc-echo /grpc-echo

CMD ["/grpc-echo"]