FROM golang:1.25.6-alpine AS builder

WORKDIR /app/src

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/auth-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

ENV GIN_MODE=release

WORKDIR /app

COPY --from=builder /usr/local/bin/auth-service /usr/local/bin/auth-service

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/auth-service"]
