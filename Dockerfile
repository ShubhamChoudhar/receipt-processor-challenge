FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o receipt-processor main.go

FROM alpine:latest

WORKDIR /root/
