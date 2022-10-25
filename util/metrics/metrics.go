package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const apiName = "telebot"

var (
	objectives   = map[float64]float64{0.5: 0.05, 0.75: .025, 0.9: 0.01, 0.95: .005, 0.99: 0.001} //nolint:gomnd
	handlerLabel = []string{"method"}

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
)
