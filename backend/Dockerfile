FROM golang:1.22 AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY internal internal
COPY main.go .
RUN go build -v -o server
RUN chmod +x server

FROM debian:bookworm-slim
COPY --from=builder /app/server /server
