package metrics

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const EnvPort = "PROMETHEUS_METRICS_PORT"
const _defaultPort = 8081

func LoadPortFromEnv() (int64, error) {
	str := os.Getenv(EnvPort)
	if str == "" {
		return _defaultPort, nil
	}
	return strconv.ParseInt(str, 10, 64)
}

func RunExporter() {
	port, err := LoadPortFromEnv()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	addr := fmt.Sprintf(":%d", port)
	// Create a new registry.
	//reg := prometheus.NewRegistry()
	// Add Go module build info.
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	//prometheus.MustRegister(collectors.NewGoCollector(
	//	collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	//))
	http.Handle("/metrics", promhttp.Handler())
	zaplogger.Sugar().Info("Hello world from new Go Collector!")
	zaplogger.Sugar().Fatal(http.ListenAndServe(addr, nil))
}
