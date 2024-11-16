package prom

import "github.com/prometheus/client_golang/prometheus"

var (
	UserCreatedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_created_count",
			Help: "Number of users created.",
		},
		[]string{"method"},
	)
	GroupCreatedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "group_created_count",
			Help: "Number of group created count.",
		},
		[]string{"method"},
	)
	DebtCreatedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "debt_created_count",
			Help: "Number of debt created count.",
		},
		[]string{"method"},
	)
)

func RegisterPrometheusMetrics() {
	prometheus.MustRegister(UserCreatedCounter, GroupCreatedCount, DebtCreatedCount)
}
