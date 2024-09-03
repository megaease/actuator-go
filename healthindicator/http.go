package healthindicator

import (
	"context"
	"net/http"
	"time"

	actuator "github.com/megaease/actuator-go/actuator"
)

// HTTPHealthIndicator checks the health of a remote HTTP endpoint.
type HTTPHealthIndicator struct {
	httpClient *http.Client
	url        string
}

// NewHTTPHealthIndicator creates a new HTTPHealthIndicator with the provided HTTP client and endpoint URL.
func NewHTTPHealthIndicator(httpClient *http.Client, url string) *HTTPHealthIndicator {
	return &HTTPHealthIndicator{
		httpClient: httpClient,
		url:        url,
	}
}

// Health checks the health of the remote HTTP endpoint and returns the health status.
func (h *HTTPHealthIndicator) Health(withDetails bool) actuator.Health {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Send a GET request to the endpoint
	requestTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return actuator.Health{
			Status: "DOWN",
			Details: &actuator.HealthDetails{
				"error": err.Error(),
			},
		}
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return actuator.Health{
			Status: "DOWN",
			Details: &actuator.HealthDetails{
				"error": err.Error(),
			},
		}
	}
	defer resp.Body.Close()

	// Determine the status based on the HTTP response
	health := actuator.Health{
		Status:  "UP",
		Details: nil,
	}
	if resp.StatusCode != http.StatusOK {
		health.Status = "DOWN"
	}

	// Include detailed information if requested
	if withDetails {
		health.Details = &actuator.HealthDetails{
			"statusCode":   resp.StatusCode,
			"status":       resp.Status,
			"url":          h.url,
			"responseTime": time.Since(requestTime).Milliseconds(),
		}
	}

	return health
}
