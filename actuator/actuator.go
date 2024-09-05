package actuator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// Actuator manages the health indicators.
type Actuator struct {
	indicators              map[string]HealthIndicator
	dynamicIndicatorsGetter func() map[string]HealthIndicator
}

// NewActuator creates a new Actuator.
func NewActuator() *Actuator {
	return &Actuator{
		indicators: make(map[string]HealthIndicator),
	}
}

// RegisterHealthIndicator registers a new health indicator.
func (a *Actuator) RegisterHealthIndicator(name string, indicator HealthIndicator) {
	a.indicators[name] = indicator
}

// RegisterDynamicHealthIndicator registers a function to dynamically get health indicators.
func (a *Actuator) RegisterDynamicHealthIndicator(dynamicFunc func() map[string]HealthIndicator) {
	a.dynamicIndicatorsGetter = dynamicFunc
}

// Health checks the health of all registered indicators with an optional `withDetails` flag.
func (a *Actuator) Health(withDetails bool) Health {
	health := Health{
		Status:     "UP",
		Components: make(map[string]Health),
	}

	for name, indicator := range a.indicators {
		componentHealth := indicator.Health(withDetails)
		health.Components[name] = componentHealth

		if componentHealth.Status == "DOWN" {
			health.Status = "DOWN"
		}
	}

	if a.dynamicIndicatorsGetter != nil {
		for name, indicator := range a.dynamicIndicatorsGetter() {
			componentHealth := indicator.Health(withDetails)
			health.Components[name] = componentHealth

			if componentHealth.Status == "DOWN" {
				health.Status = "DOWN"
			}
		}
	}

	return health
}

func checkDetailsParam(params url.Values) bool {
	// Check for both "detail" and "details"
	for _, key := range []string{"detail", "details"} {
		value, exists := params[key]
		if exists && (value[0] == "" || value[0] == "true") {
			return true
		}
	}
	return false
}

// HealthHandler returns an HTTP handler for health checks.
func (a *Actuator) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withDetails := checkDetailsParam(r.URL.Query())

		// Get the health status
		health := a.Health(withDetails)

		// Marshal the health status into a pretty-printed JSON string
		healthBytes, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to marshal health status: %v", err)))
			return
		}

		w.Write(healthBytes)
	}
}

// HealthEchoHandler returns an HTTP handler that echoes the health status.
func (a *Actuator) HealthEchoHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		withDetails := checkDetailsParam(c.QueryParams())

		// Get the health status
		health := a.Health(withDetails)

		// Marshal the health status into a pretty-printed JSON string
		healthBytes, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to marshal health status: %v", err))
		}

		// Return the health status as a plain string
		return c.String(http.StatusOK, string(healthBytes))
	}
}
