FROM golang:1.25-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app .

# Final lightweight stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /build/app .

COPY --from=builder /build/startup_funding.csv .

ENTRYPOINT ["./app"]
