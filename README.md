# API Gateway with Rate Limiting & Observability (Go)

A production-style API Gateway built in Go, designed to handle authentication, request routing, distributed rate limiting, retries, observability, and graceful shutdown for backend services.

This project focuses on real-world backend engineering concerns such as scalability, fault tolerance, reliability, and monitoring.

---

## Overview

The API Gateway acts as a single entry point for backend services.  
It centralizes cross-cutting concerns that should not be duplicated across services, such as authentication, rate limiting, retries, observability, and graceful shutdown.

Backend services remain simple and focused on business logic, while the gateway handles infrastructure-level responsibilities.

---

## Architecture

Client  
↓  
API Gateway  
├── JWT Authentication  
├── Redis-based Rate Limiting  
├── Retry with Backoff  
├── Prometheus Metrics  
↓  
Backend Services (e.g., User Service)

Public endpoints (/health, /metrics) bypass authentication.  
Protected routes enforce authentication and rate limits before proxying requests to backend services.

---

## Key Features

### Authentication
- JWT-based authentication middleware
- Token validation including signature verification and expiry checks
- User identity extracted from JWT `sub` claim
- User ID propagated safely using request context

### Rate Limiting
- Redis-backed distributed rate limiting
- Per-user request limits enforced consistently across multiple gateway instances
- Atomic request counting using Redis INCR
- Proper HTTP semantics:
  - 429 Too Many Requests
  - Retry-After response header

### Reverse Proxy
- Requests forwarded to backend services using httputil.ReverseProxy
- Gateway remains stateless and horizontally scalable

### Retry & Backoff
- Automatic retry for transient upstream failures
- Single retry with fixed backoff to prevent retry storms
- Retry count tracked using request context to avoid infinite retry loops
- Upstream failures correctly return 502 Bad Gateway

### Observability
- Prometheus metrics exposed at /metrics
- Metrics include:
  - Total HTTP requests by path, method, and status
  - Request latency histogram
  - Rate-limited request count
  - Authentication failures
  - Gateway errors
  - Retry attempts
- Metrics use low-cardinality labels to remain production-safe

### Graceful Shutdown
- Handles SIGINT and SIGTERM signals
- Stops accepting new requests on shutdown
- Allows in-flight requests to complete within a timeout
- Prevents dropped requests during deployments or restarts

---

## Tech Stack

- Language: Go
- HTTP Server: net/http
- Authentication: JWT
- Rate Limiting: Redis
- Reverse Proxy: httputil.ReverseProxy
- Observability: Prometheus
- Concurrency Model: Goroutines and request contexts
- Container Support: Docker (Redis)

---

## Request Flow

1. Client sends a request to the API Gateway  
2. Authentication middleware validates the JWT  
3. User ID is extracted and stored in the request context  
4. Rate limiter checks and increments Redis counters  
5. Request is forwarded to the backend service  
6. On transient upstream failure, the gateway retries once after backoff  
7. Metrics are recorded for latency, errors, rate limits, and retries  

---

## Running Locally

Prerequisites:
- Go 1.22 or later
- Docker (for Redis)

Steps:
1. Start Redis using Docker on port 6379
2. Start the example user service on port 8081
3. Start the API Gateway on port 8080
4. Verify using:
   - http://localhost:8080/health
   - http://localhost:8080/metrics

---

## Design Decisions

- Gateway-based authentication avoids duplicating security logic across backend services
- Redis-based rate limiting enables horizontal scalability and correctness in distributed systems
- Single retry with backoff balances resilience and system stability
- Context-based propagation avoids leaking user identity via headers
- Prometheus metrics provide observability from day one

---

## What This Project Demonstrates

- Backend system design
- Failure handling and resilience
- Distributed systems thinking
- Observability-first development
- Idiomatic Go practices
- Production-ready middleware design

---

## License

MIT
