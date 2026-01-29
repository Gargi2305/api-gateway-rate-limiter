## 🚧 Project Status

This project is currently **under active development**.

Core features like authentication, rate limiting, reverse proxying, and graceful shutdown are implemented.  
Additional improvements such as **Redis-backed distributed rate limiting, configuration cleanup, and documentation enhancements** are in progress.


# API Gateway & Distributed Rate Limiter

A backend system designed to act as a centralized API Gateway for multiple services, providing authentication, request routing, and traffic control using a Redis-backed distributed rate limiting mechanism.

## Overview
This project focuses on building a lightweight API Gateway that serves as a single entry point for backend services. It enforces security, manages request flow, and ensures fair usage through distributed rate limiting, while remaining horizontally scalable.

## Tech Stack
- Language: Python
- Framework: FastAPI
- Authentication: JWT
- Caching & Rate Limiting: Redis, Lua scripts
- Monitoring: Prometheus
- Infrastructure: Docker, NGINX

## Key Features
- JWT-based authentication middleware for request validation and user context extraction
- Distributed rate limiting using the token bucket algorithm
- Atomic quota enforcement using Redis Lua scripts
- Prometheus metrics for request throughput, latency, and rate-limit violations
- Containerized deployment behind an NGINX reverse proxy

## Architecture
The gateway sits in front of backend services, handling authentication and traffic control before forwarding requests downstream. Redis is used as a shared state store to enforce consistent rate limits across multiple gateway instances.

High-level architecture diagrams and design decisions are documented in the `/docs` directory.

## Current Status
🚧 In active development  
Repository structure and system design finalized. Core components are being implemented incrementally.

## Planned Enhancements
- Service discovery and dynamic routing
- Request logging and tracing
- Configurable per-service rate limits
- Improved error handling and retries
