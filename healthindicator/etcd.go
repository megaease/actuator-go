package healthindicator

import (
	"context"
	"time"

	actuator "github.com/megaease/actuator-go/actuator"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdHealthIndicator checks the health of an etcd cluster.
type EtcdHealthIndicator struct {
	etcdClient *clientv3.Client
}

// NewEtcdHealthIndicator creates a new EtcdHealthIndicator with the provided etcd client.
func NewEtcdHealthIndicator(etcdClient *clientv3.Client) *EtcdHealthIndicator {
	return &EtcdHealthIndicator{
		etcdClient: etcdClient,
	}
}

// Health checks the health of the etcd cluster and returns the health status.
func (e *EtcdHealthIndicator) Health(withDetails bool) actuator.Health {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initialize health response
	health := actuator.Health{
		Status:  "UP",
		Details: nil,
	}

	// Perform a status check on each etcd endpoint
	details := make(map[string]interface{})
	for _, endpoint := range e.etcdClient.Endpoints() {
		resp, err := e.etcdClient.Status(ctx, endpoint)
		if err != nil {
			health.Status = "DOWN"
			if withDetails {
				details[endpoint] = map[string]string{
					"error": err.Error(),
				}
			}
			continue
		}

		if withDetails {
			details[endpoint] = map[string]interface{}{
				"version":          resp.Version,
				"dbSize":           resp.DbSize,
				"leader":           resp.Leader,
				"raftIndex":        resp.RaftIndex,
				"raftTerm":         resp.RaftTerm,
				"isLearner":        resp.IsLearner,
				"memberId":         resp.Header.MemberId,
				"raftAppliedIndex": resp.RaftAppliedIndex,
			}
		}
	}

	if withDetails {
		health.Details = &actuator.HealthDetails{
			"endpoints": details,
		}
	}

	return health
}
