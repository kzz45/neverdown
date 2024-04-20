package env

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s.io/klog/v2"
)

func getHostName() (string, error) {
	cmd := exec.Command("/bin/hostname")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	hostname := strings.TrimSpace(string(out))
	if hostname == "" {
		return "", fmt.Errorf("no hostname get from cmd '/bin/hostname' in the container, please check")
	}
	return hostname, nil
}

func GetHostNameMustSpecified() string {
	t, err := getHostName()
	if err != nil {
		klog.Fatal(err)
	}
	return t
}

func GetHostName() (string, error) {
	hostName := os.Getenv("HOST_NAME")
	if hostName != "" {
		return hostName, nil
	}
	klog.Info("get HOST_NAME from env failed, is env.(\"HOST_NAME\") already set? Will use hostname instead")
	return getHostName()
}
