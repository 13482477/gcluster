package metric

import (
	"context"
	"time"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/endpoint"
)

type GClusterMetric struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
}

func MCloudMetricServer(method string, metric *GClusterMetric) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if metric != nil {
				defer func(begin time.Time) {
					metric.RequestLatency.With("method", method).Observe(time.Since(begin).Seconds())
				}(time.Now())
				metric.RequestCount.With("method", method).Add(1)
			}
			return next(ctx, request)
		}
	}
}
