FROM golang:1.21 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o auth_usecase .

FROM alpine:latest

RUN mkdir /app_repository

COPY --from=builder /app/auth /app

WORKDIR /app

EXPOSE 7002

CMD ./auth