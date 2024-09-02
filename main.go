package main

import (
	"log"
	"net/http"

	"github.com/megaease/actuator-go/actuator"
	"github.com/megaease/actuator-go/healthindicator"
)

func main() {
	actuator := actuator.NewActuator()

	// Register indicators
	actuator.RegisterHealthIndicator("memory", &healthindicator.MemoryHealthIndicator{})

	// Create HTTP server
	http.HandleFunc("/actuator/health", actuator.HealthHandler())

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
