# Dockerfile for hotels-api
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o hotels .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/hotels .
EXPOSE 8080
ENTRYPOINT [ "/app/hotels" ]
