package prometheus

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type RequestsToNodeCount struct {
	internal *prometheus.CounterVec
}

func NewRequestsToNodeCount(namespace string) *RequestsToNodeCount {
	return &RequestsToNodeCount{
		internal: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "node_request_count",
		}, []string{
			"chain",
		}),
	}
}

func (r *RequestsToNodeCount) Inc(chain entities.Chain) {
	r.internal.WithLabelValues(string(chain)).Inc()
}
