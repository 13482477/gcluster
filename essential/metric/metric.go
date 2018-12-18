package metric

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/endpoint"
	"context"
	"time"
)

type MCloudMetric struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
}

func MCloudMetricServer(method string, metric *MCloudMetric) endpoint.Middleware {
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
