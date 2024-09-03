package healthindicator

import actuator "github.com/megaease/actuator-go/actuator"

// SelfHealthIndicator is a simple health indicator that always returns "UP".
type SelfHealthIndicator struct{}

// NewSelfHealthIndicator creates a new SelfHealthIndicator.
func NewSelfHealthIndicator() *SelfHealthIndicator {
	return &SelfHealthIndicator{}
}

// Health always returns an "UP" health status.
func (s *SelfHealthIndicator) Health(withDetails bool) actuator.Health {
	health := actuator.Health{
		Status: "UP",
	}

	// If withDetails is true, include an empty details structure
	if withDetails {
		health.Details = &actuator.HealthDetails{}
	}

	return health
}
