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
- Name
- Email

Do NOT implement any functionality that is not listed above.

---------------------------------------------------------------------

# Service Responsibilities

Auth Service is responsible only for:

- user registration
- user login
- user logout
- storing user account data:
  - name
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
- name
- email

If password storage is required for authentication, it may be used internally,
but it must not expand the domain scope beyond the roadmap.


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