FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY microservices/wiki-service ./microservices/wiki-service

WORKDIR /app/microservices/wiki-service
ENV GOFLAGS=-mod=mod
ENV GOWORK=off
RUN go build -o /wiki-service ./cmd/main.go

FROM alpine:latest
COPY --from=builder /wiki-service /wiki-service
EXPOSE 50051 8080
CMD ["/wiki-service"]