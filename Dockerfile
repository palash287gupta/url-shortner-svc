# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortner-svc .

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortner-svc .

EXPOSE 8080

ENTRYPOINT ["./url-shortner-svc"]
