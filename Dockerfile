FROM golang:1.20-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN CGO