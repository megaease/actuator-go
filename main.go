package main

import (
	"log"
	"net/http"

	"github.com/megaease/actuator-go/actuator"
	"github.com/megaease/actuator-go/healthindicator"
)

func main() {
	actuatorExample := actuator.NewActuator()

	// Register indicators
	actuatorExample.RegisterHealthIndicator("self-example", &healthindicator.SelfHealthIndicator{})

	actuatorExample.RegisterDynamicHealthIndicator(func() map[string]actuator.HealthIndicator {
		return map[string]actuator.HealthIndicator{
			"dynamic-indicator-1": &healthindicator.SelfHealthIndicator{},
		}
	})

	// Create HTTP server
	http.HandleFunc("/actuator/health", actuatorExample.HealthHandler())

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
