package healthindicator

import (
	"runtime"

	actuator "github.com/megaease/actuator-go/actuator"
)

// MemoryHealthIndicator checks the system's memory usage.
type MemoryHealthIndicator struct{}

// Health checks the memory usage and returns the appropriate health status.
func (m *MemoryHealthIndicator) Health(withDetails bool) actuator.Health {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	health := actuator.Health{
		Status:     "UP",
		Components: nil,
		Details:    nil,
	}

	if withDetails {
		health.Details = &actuator.HealthDetails{
			"Alloc":      memStats.Alloc,
			"TotalAlloc": memStats.TotalAlloc,
			"Sys":        memStats.Sys,
			"NumGC":      memStats.NumGC,
		}
	}

	return health
}
