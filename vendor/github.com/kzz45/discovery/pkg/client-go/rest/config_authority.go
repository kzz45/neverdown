package rest

import (
	"errors"
	"net"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"
)

var ErrAuthorityNotInCluster = errors.New("unable to load in-cluster configuration, AUTHORITY_SERVICE_HOST and AUTHORITY_SERVICE_PORT must be defined")
var ErrAuthorityCAFileNotInCluster = errors.New("unable to load in-cluster configuration, AUTHORITY_SERVICE_CAFILE must be defined")

func InDicoveryClusterConfig() (*Config, error) {
	const (
		tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	// rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("DISCOVERY_SERVICE_HOST"), os.Getenv("DISCOVERY_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, ErrNotInCluster
	}

	// token, err := os.ReadFile(tokenFile)
	// if err != nil {
	// 	return nil, err
	// }
	rootCAFile := os.Getenv("DISCOVERY_SERVICE_CAFILE")
	if rootCAFile == "" {
		return nil, ErrCAFileNotInCluster
	}

	tlsClientConfig := TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		ContentConfig: ContentConfig{
			ContentType: runtime.ContentTypeProtobuf,
		},
		// BearerToken:     string(token),
		// BearerTokenFile: tokenFile,
	}, nil
}

func InAuthorityClusterConfig() (*Config, error) {
	const (
	// tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	//rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("AUTHORITY_SERVICE_HOST"), os.Getenv("AUTHORITY_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, ErrAuthorityNotInCluster
	}

	rootCAFile := os.Getenv("AUTHORITY_SERVICE_CAFILE")
	if rootCAFile == "" {
		return nil, ErrAuthorityCAFileNotInCluster
	}

	//token, err := ioutil.ReadFile(tokenFile)
	//if err != nil {
	//	return nil, err
	//}

	tlsClientConfig := TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Fatal("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		ContentConfig: ContentConfig{
			ContentType: runtime.ContentTypeProtobuf,
		},
		//BearerToken:     string(token),
		//BearerTokenFile: tokenFile,
	}, nil
}
