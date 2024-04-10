package env

import (
	"os"
	"strconv"

	"k8s.io/klog/v2"
)

const (
	ExposePProfDebug = "PPROF_DEBUG"
)

// PProfDebug determines that whether the server will expose pprof routers
func PProfDebug() bool {
	a := os.Getenv(ExposePProfDebug)
	if a == "" {
		return false
	}
	t, err := strconv.ParseBool(a)
	if err != nil {
		klog.Error(err)
		return false
	}
	return t
}
