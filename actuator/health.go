package actuator

// HealthDetails represents additional information about the health status.
type HealthDetails map[string]interface{}

// Health represents the health status of a component.
type Health struct {
	Status     string            `json:"status"`
	Components map[string]Health `json:"components,omitempty"`
	Details    *HealthDetails    `json:"details,omitempty"`
	Groups     []string          `json:"groups,omitempty"`
}

// HealthIndicator is the interface that should be implemented by any health check.
type HealthIndicator interface {
	Health(withDetails bool) Health
}
