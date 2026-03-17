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
* storing user account data

The service maintains basic user identity information and provides authentication capabilities for the platform.

---

## Current Scope (Roadmap Stage: Base)

The following functionality is implemented in the current stage.

### User Account

Supported operations:

* **Registration**
* **Login**
* **Logout**

User data stored by the service:

* `id`
* `username`
* `email`

---

## Local Run

1. Create a local environment file:

```bash
cp .example.env .env
```

2. Start PostgreSQL in Docker:

```bash
docker run -d \
  --name auth-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=auth_service \
  -p 5432:5432 \
  postgres:16
```

3. Initialize the required tables:

```bash
docker exec -i auth-postgres psql -U postgres -d auth_service -c "
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
"
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
curl http://localhost:8080/
```

The current local entrypoint uses temporary stubs for the logger and password hasher.

---

## TODO

* implement password validation rules in `ValidatePassword`
* replace temporary stubs in `src/internal/logger.go` and `src/internal/adapters/hasher.go` with real implementations during application wiring
