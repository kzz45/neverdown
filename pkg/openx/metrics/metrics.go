package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	_namespace = "openx_controller"
)

const (
	_labelNamespace        = "namespace"
	_labelGroupVersionKind = "group_version_kind"
	_labelVerb             = "verb"
)

type Collector interface {
	IncrOnlinePlayer()
	DecOnlinePlayer()
	AccumulateRequest(namespace string, groupVersionKind string, verb string)
}

type collector struct {
	onlinePlayer prometheus.Gauge
	requests     *prometheus.CounterVec
}

func NewPrometheusCollector() Collector {
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	if podNamespace == "" {
		podNamespace = "kube-neverdown"
	}
	constLabels := prometheus.Labels{
		"controller_pod":       podName,
		"controller_namespace": podNamespace,
	}
	collector := &collector{
		onlinePlayer: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   _namespace,
			Name:        "online_player",
			Help:        "Number of the online players",
			ConstLabels: constLabels,
		}),
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   _namespace,
			Name:        "requests",
			Help:        "Number of the requests to openx",
			ConstLabels: constLabels,
		}, []string{_labelNamespace, _labelGroupVersionKind, _labelVerb},
		),
	}
	// Since we use the DefaultRegisterer, in test cases, the metrics
	// might be registered duplicately, unregister them before re register.
	prometheus.Unregister(collector.onlinePlayer)
	prometheus.Unregister(collector.requests)
	prometheus.MustRegister(
		collector.onlinePlayer,
		collector.requests,
	)
	return collector
}

func (c *collector) IncrOnlinePlayer() {
	c.onlinePlayer.Inc()
}

func (c *collector) DecOnlinePlayer() {
	c.onlinePlayer.Dec()
}

func (c *collector) AccumulateRequest(namespace string, groupVersionKind string, verb string) {
	c.requests.With(prometheus.Labels{
		_labelNamespace:        namespace,
		_labelGroupVersionKind: groupVersionKind,
		_labelVerb:             verb,
	}).Inc()
}
