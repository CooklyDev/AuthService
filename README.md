# Auth Service

## Overview

Auth Service is a microservice responsible for user account authentication in the Cookly platform.

The service manages the authentication lifecycle of a user account and handles user registration and login.

## Run

1. Create a local environment file:

```bash
cp .example.env .env
```

2. Start PostgreSQL, run migrations, and start Auth Service:

```bash
docker compose up --build
```

By default, the service will be available at `http://localhost:8080`.

Swagger UI can be accessed at `http://localhost:8080/swagger/index.html`.

## Development

### Generate OpenAPI files:

```bash
cd src
./generate_swagger.sh
```

This does not require a globally installed `swag` binary.

### Run tests:

```bash
go test -v ./...
```

### Run linter and formatter:

```bash
go fmt ./...
golangci-lint run
```
