# API Gateway & Rate Limiter

A backend system built in Go that acts as a centralized API Gateway for backend services, providing authentication, request routing, rate limiting, and graceful shutdown.

## 🚧 Project Status

This project is currently **under active development**.

Core features such as authentication, request forwarding, in-memory rate limiting, and graceful shutdown are implemented.  
Enhancements like **Redis-backed distributed rate limiting, configuration cleanup, and observability** are planned.

## Overview

This project focuses on building a lightweight API Gateway that serves as a single entry point for backend services.  
It enforces security, manages request flow, and ensures fair usage through per-user rate limiting, while being designed for horizontal scalability.

## Tech Stack

- Language: **Go**
- HTTP Framework: **net/http**
- Authentication: **JWT**
- Rate Limiting: **In-memory (Redis-based planned)**
- Reverse Proxy: **httputil.ReverseProxy**
- Concurrency & Safety: **Go routines, sync primitives**

## Key Features

- JWT-based authentication middleware for request validation
- Secure user identity propagation using request context
- Per-user rate limiting with proper HTTP semantics (429, Retry-After)
- Reverse proxying of requests to upstream services
- Graceful shutdown using OS signals (SIGINT, SIGTERM)
- Clean separation between gateway and downstream services

## Architecture

Client → API Gateway → Backend Services (e.g., User Service)

The API Gateway handles authentication, rate limiting, and request routing before forwarding requests to downstream services.  
User identity is propagated securely using request context instead of headers.

## Planned Enhancements

- Redis-backed distributed rate limiting
- External configuration management
- Metrics and observability (Prometheus)
- Dockerized deployment
- Improved logging and tracing
