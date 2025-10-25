<<<<<<< HEAD
# E-commerce Real-time Monitoring Service

A simplified demo of a real-time monitoring service for e-commerce orders with WebSocket streaming, Redis pub/sub, and Prometheus metrics.

## Features

✅ **Gorilla WebSocket** - Real-time bidirectional communication
✅ **Redis Pub/Sub** - Message queuing and distribution
✅ **Prometheus Metrics** - Queue depth, latency, error rates
✅ **Grafana Dashboard** - Real-time visualization
✅ **Concurrent Connections** - Handles 200+ WebSocket connections
✅ **Goroutines** - Efficient concurrent processing

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   E-commerce    │───▶│   Redis Pub/Sub │───▶│  WebSocket Hub  │
│   Order Events  │    │   (Scaling)     │    │  (200+ conns)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Prometheus    │    │   Dashboard     │
                       │   (Metrics)     │    │   (Real-time)   │
                       └─────────────────┘    └─────────────────┘
```

## Quick Start

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Start Redis** (optional - demo works without it)
   ```bash
   docker run -d -p 6379:6379 redis:alpine
   ```

3. **Run the Service**
   ```bash
   go run main.go
   ```

4. **Open Dashboard**
   - Dashboard: http://localhost:8080
   - Metrics: http://localhost:8080/metrics
   - WebSocket: ws://localhost:8080/ws

## Key Components

### WebSocket Hub
- Manages 200+ concurrent connections
- Thread-safe connection handling
- Automatic cleanup and reconnection

### Redis Pub/Sub
- Horizontal scaling across instances
- Message queuing and distribution
- Fault tolerance and reliability

### Prometheus Metrics
- `orders_total` - Total orders by status
- `websocket_connections_active` - Active connections
- `order_processing_latency_seconds` - Processing latency

### Real-time Dashboard
- Live order statistics
- Revenue tracking
- Error rate monitoring
- Queue depth visualization

## Production Considerations

This is a **simplified demo**. For production use, consider:

- **Authentication** - JWT tokens, API keys
- **Rate Limiting** - Prevent abuse
- **Error Handling** - Graceful degradation
- **Monitoring** - Health checks, alerts
- **Security** - HTTPS, CORS, input validation
- **Scaling** - Load balancers, multiple instances
- **Persistence** - Database integration
- **Testing** - Unit tests, integration tests

## Demo Limitations

- Simulated data (not real e-commerce integration)
- No authentication or security
- Basic error handling
- Simplified Redis configuration
- Mock metrics generation

## Next Steps

To make this production-ready:

1. **Integrate with real e-commerce system**
2. **Add authentication and authorization**
3. **Implement proper error handling**
4. **Add comprehensive testing**
5. **Set up monitoring and alerting**
6. **Configure for horizontal scaling**
=======
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
>>>>>>> 7c8d21818a1d9dbf7170d721c2e8b5d5fc9b22e8
