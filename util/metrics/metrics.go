package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	apiName = "telebot"
)

var (
	objectives     = map[float64]float64{0.5: 0.05, 0.75: .025, 0.9: 0.01, 0.95: .005, 0.99: 0.001} //nolint:gomnd
	expenseBuckets = []float64{1, 10, 100, 1000, 5000, 10000, 100000}                               //nolint:gomnd
	handlerLabel   = []string{"method"}

	RequestOpsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: apiName,
		Name:      "requests_total",
	}, handlerLabel)

	RequestDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  apiName,
		Name:       "requests_duration_seconds",
		Objectives: objectives,
	}, handlerLabel)

	RequestOpsInternalError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: apiName,
		Name:      "requests_error_total",
	}, handlerLabel)

	ExpenseCounter = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: apiName,
		Name:      "expenses_amount",
		Buckets:   expenseBuckets,
	})

	WrongUsageCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: apiName,
		Name:      "wrong_usage_total",
	})

	CacheMissCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: apiName,
		Name:      "cache_miss_counter",
	})

	CacheHitCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: apiName,
		Name:      "cache_hit_counter",
	})
)
