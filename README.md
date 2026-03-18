# Auth Service

## Overview

Auth Service is a microservice responsible for user account authentication in the Cookly platform.

Cookly is a web platform that allows users to manage cooking recipes and share them with other users.

Auth Service manages the authentication lifecycle of a user account and handles user registration and login.

---

## Responsibilities

Auth Service is responsible for:

* user registration
* user login
* user logout
* storing user account data

The service maintains basic user identity information and provides authentication capabilities for the platform.

---

## Local Run

1. Create a local environment file:

```bash
cp .example.env .env
```

`APP_PORT` controls the HTTP port for the Auth Service. The default value is `8080`.

2. Start PostgreSQL only:

```bash
docker compose up -d postgres
```

3. Apply migrations:

```bash
docker compose run --rm migrate
```

4. Load environment variables in the current shell from `src`:

```bash
cd src
source ./set_env.sh
```

5. Run the service:

```bash
go run ./cmd
```

6. Check that the service is reachable:

```bash
curl http://localhost:${APP_PORT}/
```

The current local entrypoint uses `ConsoleLogger` and bcrypt password hashing.

---

## Docker Compose

1. Create a local environment file:

```bash
cp .example.env .env
```

2. Start PostgreSQL, run migrations, and start Auth Service:

```bash
docker compose up --build
```

3. Check that the service is reachable:

```bash
curl http://localhost:${APP_PORT:-8080}/
```

4. Stop the stack:

```bash
docker compose down
```

If you need to remove the PostgreSQL volume too:

```bash
docker compose down -v
```

---

## Migrations

Apply migrations:

```bash
docker compose run --rm migrate
```

Check migration status:

```bash
docker compose run --rm migrate status
```

Rollback the last migration:

```bash
docker compose run --rm migrate down
```

---

## Swagger Generation

1. Generate OpenAPI files:

```bash
cd src
swag init -g cmd/main.go -o docs
```

2. Access the Swagger UI at:

```text
http://localhost:${APP_PORT:-8080}/swagger/index.html
```
