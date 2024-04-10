package env

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kzz45/neverdown/pkg/zaplogger"
)

const (
	NodeName          = "NODE_NAME"
	HostIP            = "HOST_IP"
	PodName           = "POD_NAME"
	PodNamespace      = "POD_NAMESPACE"
	PodIP             = "POD_IP"
	PodServiceAccount = "POD_SERVICE_ACCOUNT"
)

// getHostName gets the hostname of the host machine if the container is started by docker run --net=host
func getHostName() (string, error) {
	cmd := exec.Command("/bin/hostname")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	hostname := strings.TrimSpace(string(out))
	if hostname == "" {
		return "", fmt.Errorf("no hostname get from cmd '/bin/hostname' in the container, please check")
	}
	return hostname, nil
}

// GetHostName gets the hostname of host machine
func GetHostName() (string, error) {
	hostName := os.Getenv("HOST_NAME")
	if hostName != "" {
		return hostName, nil
	}
	zaplogger.Sugar().Info("get HOST_NAME from env failed, is env.(\"HOST_NAME\") already set? Will use hostname instead")
	return getHostName()
}

// GetHostNameMustSpecified will fatal if the hostname hasn't been specified
func GetHostNameMustSpecified() string {
	t, err := getHostName()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return t
}

// GetPodIP returns the ip which has been allocated to the pod in the k8s cluster
func GetPodIP() (string, error) {
	podIp := os.Getenv(PodIP)
	if podIp != "" {
		return podIp, nil
	}
	return "", fmt.Errorf("no env variable getting from '%s' in the container", PodIP)
}

// GetPodNamespace returns the namespace which has been allocated to the pod in the k8s cluster
func GetPodNamespace() (string, error) {
	namespace := os.Getenv(PodNamespace)
	if namespace != "" {
		return namespace, nil
	}
	return "default", fmt.Errorf("no env variable getting from '%s' in the container", PodNamespace)
}

// GetPodNamespaceMustSpecified will fatal if the namespace hasn't been specified
func GetPodNamespaceMustSpecified() string {
	t, err := GetPodNamespace()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return t
}
