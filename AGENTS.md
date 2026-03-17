# Project Context

Project: Cookly

Cookly is a web platform for managing cooking recipes and sharing them with other users.

The system is built using a microservice architecture.

Currently the system contains two active services in development:
- Auth Service
- Content Service

Each service owns its own database.

External communication between services and clients is done through REST APIs.
Internal service-to-service communication may use gRPC.


---------------------------------------------------------------------

# Service Context

Service: Auth Service

Auth Service is responsible only for user account authentication.

This service manages user registration, login and logout.

The current development scope is strictly limited to what is defined in the roadmap.


---------------------------------------------------------------------

# Roadmap Scope (DO NOT EXTEND)

Base → User → Account:

- Registration
- Login / Logout
- Username
- Email

Do NOT implement any functionality that is not listed above.

---------------------------------------------------------------------

# Service Responsibilities

Auth Service is responsible only for:

- user registration
- user login
- user logout
- storing user account data:
  - username
  - email

The service should not handle:

- recipes
- collections
- ingredients
- social features
- discovery features
- media


---------------------------------------------------------------------

# Minimum Domain Model

User entity must contain only:

- id
- username
- email
- password



---------------------------------------------------------------------

# Technology Stack

Backend:
- Go

Database:
- PostgreSQL

Cache / session storage:
- Redis


---------------------------------------------------------------------

# Architecture

The project must follow Clean Architecture principles.

Dependencies must always point inward.

Outer layers must depend on inner layers, never the opposite.

The code must be divided into the following layers.


Entities
Enterprise-level business objects.
Contain core business rules.
Must not depend on frameworks, databases, or external libraries.


Use Cases
Implement system behavior and application rules.
Operate on entities.
Define interfaces for external dependencies.


Adapters
Implement the interfaces defined by use cases.

Examples:
- repositories
- database adapters
- HTTP adapters
- external service clients

Adapters must not implement business logic or any behavior outside the contract defined by the interface they satisfy.


Frameworks
External delivery mechanisms.

Examples:
- HTTP handlers
- gRPC handlers
- consumers
- web frameworks
- routing


Main
Application assembly layer.

Responsibilities:
- wiring dependencies
- initializing infrastructure
- starting the application

Main is the lowest-level policy and entry point into the system.


---------------------------------------------------------------------

# Suggested Project Structure

... not defined yet

---------------------------------------------------------------------

# Coding Guidelines

Prefer small functions with one clear responsibility.

Avoid deeply nested conditionals.

Reuse existing abstractions before creating new ones.

Keep control flow easy to scan.

Avoid overly dense expressions.

Prefer explicit code over clever code.

Readable code is more important than clever optimizations.


---------------------------------------------------------------------

# Change Policy

Prefer minimal, reviewable diffs.

Do not refactor unrelated code.

Do not change public APIs unless explicitly requested.

Do not introduce new dependencies unless necessary.

Changes must stay within the defined service scope.


---------------------------------------------------------------------

# Validation Rules

Before finishing any change:

- run code formatting
- run unit tests for the affected module
- run integration tests if database behavior changed

The task must not be considered complete if tests fail.


---------------------------------------------------------------------

# Command Usage

Use commands from the correct directory:

- run `pre-commit` commands from the repository root
- run Go module commands from `src`

Prefer the following commands:

- formatting: `cd src && gofmt -w <files>`
- unit tests: `cd src && go test ./...`
- linter: `cd src && golangci-lint run -c ../.golangci.yml ./...`
- pre-commit full check: `pre-commit run --all-files`

Do not replace these commands with alternative tools unless explicitly requested.


---------------------------------------------------------------------

# Logger Usage

Use the logger in use cases, adapters, and application startup code.

Do not log from domain entities.

Use log levels consistently:

- `Debug` for development diagnostics and branch-level execution details
- `Info` for successful service lifecycle events and successful business operations
- `Warn` for expected but undesirable situations such as invalid input, duplicate email, or failed login
- `Error` for technical failures that prevent completing an operation

Log only useful operational context:

- operation name such as `register`, `login`, `logout`
- `user_id` when available
- masked `email` when needed
- dependency name such as `postgres` or `redis`
- request correlation identifiers when available
- sanitized error details

Never log sensitive data:

- passwords
- password hashes
- access tokens
- refresh tokens
- session identifiers in plain form
- secrets
- full request bodies

Prefer structured fields over string concatenation.


---------------------------------------------------------------------

# Testing Rules

Business logic correctness is the highest priority.

Implementation should follow TDD where feasible.


--------------------------------------------------

Unit Tests

Unit tests must:

- test a small isolated piece of logic
- run quickly
- avoid external dependencies

Use mocks and stubs for external dependencies.

For use case unit tests:

- prefer tests for business branches and observable outcomes
- do not test whether mocks or stubs were called internally
- treat mocks and stubs as dumb test helpers by default
- do not put business logic or meaningful branching into mocks or stubs
- do not add separate tests for mocks or stubs
- do not add tests that only verify passthrough of artificially injected dependency errors unless that behavior is itself the business rule being requested
- keep stubs dumb: they should return predefined data and should not contain callback-driven logic
- do not assert internal state captured by stubs when the same behavior can be verified through the use case result

Follow the AAA pattern:

Arrange  
Prepare dependencies and input data.

Act  
Execute the tested unit.

Assert  
Verify the results.

Comments marking AAA stages must be present in the test code.

Tests should verify behavior rather than implementation details.


--------------------------------------------------

Integration Tests

Integration tests may interact with real infrastructure such as the database.

Each integration test must:

- restore the database to a clean state
- not depend on the order of execution


---------------------------------------------------------------------

# AI Behavior Rules

When generating code for this service:

- stay strictly within the roadmap scope
- do not add new features
- do not expand the domain model
- follow Clean Architecture layering
- produce minimal, readable implementations
- focus on working MVP functionality

If a design decision is unclear, prefer the simplest implementation that satisfies the roadmap.

This file defines the AI development context for the repository.
All AI coding agents must follow these rules when generating or modifying code.
