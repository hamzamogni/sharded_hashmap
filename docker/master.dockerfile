FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/master /app/cmd/master

FROM alpine:3.12

COPY --from=builder /app/master /app/master

EXPOSE 8080
CMD ["/app/master"]

