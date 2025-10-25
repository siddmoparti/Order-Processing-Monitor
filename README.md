# Real-Time Order Processing Monitor

A production-grade real-time monitoring system for e-commerce order pipelines — built with **Go**, **Redis**, **Gorilla WebSocket**, **Docker**, and **Prometheus + Grafana** for observability.

This service streams **live order events** to a browser dashboard with **millisecond latency**, enabling instant visibility into queue health, latency trends, and failure conditions.

---

## 🚀 Features

- **Real-time WebSocket dashboard** using Go + Gorilla WebSocket
- **Redis Pub/Sub ingestion layer** decoupling producers from WebSocket broadcast
- **Concurrency-safe fanout with goroutines** and graceful connection lifecycle
- **Live metrics tracking** — queue depth, p95 latency, failure rates
- **Prometheus + Grafana integration** for production-grade observability
- **Dockerized deployment** with `docker-compose up --build`

---

## 🧩 Architecture

```mermaid
flowchart LR
    Producer -->|publishes| Redis[(Redis Pub/Sub)]
    Redis -->|stream| GoService[Go WebSocket Service]
    GoService -->|WS push| Dashboard[Web Client]
    GoService -->|metrics| Prometheus
    Prometheus --> Grafana
