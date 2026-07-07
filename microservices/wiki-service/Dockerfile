FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.work go.work.sum ./
COPY microservices/wiki-service ./microservices/wiki-service
RUN go build -o /wiki-service ./microservices/wiki-service/cmd/main.go

FROM alpine:latest
COPY --from=builder /wiki-service /wiki-service
EXPOSE 50051 8080
CMD ["/wiki-service"]