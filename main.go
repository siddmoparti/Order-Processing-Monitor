package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/go-redis/redis/v8"
	"context"
)

// Order represents an e-commerce order
type Order struct {
	ID        string    `json:"id"`
	Customer  string    `json:"customer"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Stats represents real-time statistics
type Stats struct {
	TotalOrders    int     `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	ActiveOrders   int     `json:"active_orders"`
	AverageOrder   float64 `json:"average_order"`
	ErrorRate      float64 `json:"error_rate"`
	QueueDepth     int     `json:"queue_depth"`
}

// WebSocket connection manager
type Hub struct {
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan []byte
	mu         sync.RWMutex
	redis      *redis.Client
}

// Prometheus metrics
var (
	ordersTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Total number of orders processed",
		},
		[]string{"status"},
	)

	websocketConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Number of active WebSocket connections",
		},
	)

	orderLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "order_processing_latency_seconds",
			Help: "Order processing latency",
		},
	)
)

func init() {
	prometheus.MustRegister(ordersTotal)
	prometheus.MustRegister(websocketConnections)
	prometheus.MustRegister(orderLatency)
}

func newHub() *Hub {
	// Initialize Redis client (simplified - in real app you'd configure properly)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan []byte),
		redis:      rdb,
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			websocketConnections.Inc()
			log.Printf("Client connected. Total connections: %d", len(h.clients))

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
			websocketConnections.Dec()
			log.Printf("Client disconnected. Total connections: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for conn := range h.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Simulate order processing with Redis pub/sub
func (h *Hub) processOrders() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Simulate new order
			order := Order{
				ID:        fmt.Sprintf("order_%d", time.Now().Unix()),
				Customer: fmt.Sprintf("customer_%d", rand.Intn(100)),
				Amount:    rand.Float64() * 1000,
				Status:    []string{"pending", "processing", "completed", "failed"}[rand.Intn(4)],
				Timestamp: time.Now(),
			}

			// Publish to Redis (simplified)
			orderJSON, _ := json.Marshal(order)
			h.redis.Publish(context.Background(), "orders", orderJSON)

			// Update metrics
			ordersTotal.WithLabelValues(order.Status).Inc()

			// Simulate processing latency
			latency := time.Duration(rand.Intn(1000)) * time.Millisecond
			orderLatency.Observe(latency.Seconds())

			// Generate stats and broadcast
			stats := h.generateStats()
			statsJSON, _ := json.Marshal(stats)
			h.broadcast <- statsJSON
		}
	}
}

func (h *Hub) generateStats() Stats {
	// Simplified stats generation
	return Stats{
		TotalOrders:  rand.Intn(1000) + 500,
		TotalRevenue: rand.Float64() * 100000,
		ActiveOrders: rand.Intn(50) + 10,
		AverageOrder: rand.Float64() * 200,
		ErrorRate:    rand.Float64() * 0.05,
		QueueDepth:   rand.Intn(20) + 5,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	hub.register <- conn

	// Keep connection alive
	go func() {
		defer func() {
			hub.unregister <- conn
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}
		}
	}()
}

func main() {
	hub := newHub()
	go hub.run()
	go hub.processOrders()

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Simple dashboard endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>E-commerce Monitoring Dashboard</title>
    <script>
        const ws = new WebSocket('ws://localhost:8080/ws');
        ws.onmessage = function(event) {
            const stats = JSON.parse(event.data);
            document.getElementById('total-orders').textContent = stats.total_orders;
            document.getElementById('total-revenue').textContent = '$' + stats.total_revenue.toFixed(2);
            document.getElementById('active-orders').textContent = stats.active_orders;
            document.getElementById('average-order').textContent = '$' + stats.average_order.toFixed(2);
            document.getElementById('error-rate').textContent = (stats.error_rate * 100).toFixed(2) + '%';
            document.getElementById('queue-depth').textContent = stats.queue_depth;
        };
    </script>
</head>
<body>
    <h1>E-commerce Monitoring Dashboard</h1>
    <div>
        <h2>Real-time Stats</h2>
        <p>Total Orders: <span id="total-orders">0</span></p>
        <p>Total Revenue: <span id="total-revenue">$0</span></p>
        <p>Active Orders: <span id="active-orders">0</span></p>
        <p>Average Order: <span id="average-order">$0</span></p>
        <p>Error Rate: <span id="error-rate">0%</span></p>
        <p>Queue Depth: <span id="queue-depth">0</span></p>
    </div>
    <p><a href="/metrics">Prometheus Metrics</a></p>
</body>
</html>
		`)
	})

	log.Println("Starting server on :8080")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws")
	log.Println("Dashboard: http://localhost:8080")
	log.Println("Metrics: http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
