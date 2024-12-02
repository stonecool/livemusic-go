package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ChromeStateGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "chrome_state",
			Help: "Current state of chrome instances",
		},
		[]string{"id", "addr"},
	)

	ChromeTaskCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chrome_tasks_total",
			Help: "Total number of tasks executed by chrome instances",
		},
		[]string{"id", "category"},
	)
)

func init() {
	prometheus.MustRegister(ChromeStateGauge)
	prometheus.MustRegister(ChromeTaskCounter)
}
