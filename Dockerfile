# build
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Then copy the rest
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["/app/main"]
