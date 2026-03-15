# Build Stage
FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/url-shortener/main.go

# Runtime Stage
FROM scratch
COPY --from=builder /app/main .
EXPOSE 8082
ENTRYPOINT ["./main"]