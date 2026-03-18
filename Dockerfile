FROM golang:1.25.6-alpine AS builder

WORKDIR /app/src

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/auth-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

ARG APP_PORT=8080
ENV GIN_MODE=release
ENV APP_PORT=${APP_PORT}

WORKDIR /app

COPY --from=builder /usr/local/bin/auth-service /usr/local/bin/auth-service

USER appuser

EXPOSE ${APP_PORT}

ENTRYPOINT ["/usr/local/bin/auth-service"]
