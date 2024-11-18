# Dockerfile for hotels-api
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o search .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/search .
EXPOSE 8082
ENTRYPOINT [ "/app/search" ]
