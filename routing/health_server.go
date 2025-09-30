package routing

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// livenessHandler responds with 200 OK if the application is running.
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// readinessHandler responds with 200 OK if the application is ready to serve traffic.
// You can add more sophisticated checks here later (e.g., database connectivity).
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// For now, just a simple OK.
	// TODO: Add checks for dependencies (DB, external services, etc.)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// StartHealthServer starts an HTTP server for health and readiness probes.
// It's intended to be run in a goroutine.
func StartHealthServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/livez", livenessHandler)
	mux.HandleFunc("/readyz", readinessHandler)
	
	// Add Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	serverAddr := ":" + port
	log.Printf("Starting health and readiness server on %s", serverAddr)
	log.Printf("Prometheus metrics available at http://localhost%s/metrics", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Health server failed: %v", err)
	}
}
