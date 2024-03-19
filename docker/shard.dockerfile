FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/shard /app/cmd/shard

FROM alpine:3.12

COPY --from=builder /app/shard /shard

EXPOSE 9090

CMD ["/shard"]
