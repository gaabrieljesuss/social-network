package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Número total de requisições recebidas",
		},
		[]string{"method", "endpoint"},
	)

	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_response_time_seconds",
			Help:    "Tempo de resposta das requisições",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func Init() {
	fmt.Println("Registrando métricas no Prometheus...")
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(ResponseTime)
}
