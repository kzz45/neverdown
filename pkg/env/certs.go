package env

import (
	"crypto/tls"
	"errors"
	"os"
)

const (
	StaticCertFile   = "TLS_OPTION_CERT_FILE"
	StaticKeyFile    = "TLS_OPTION_KEY_FILE"
	StaticCaFile     = "TLS_OPTION_CA_FILE"
	StaticServerName = "TLS_OPTION_SERVER_NAME"
)

var ErrCAFileNotInCluster = errors.New("unable to load tls configuration, TLS_OPTION_CERT_FILE or TLS_OPTION_KEY_FILE must be defined")
var ErrGrpcTLSNotInCluster = errors.New("unable to load tls configuration, TLS_OPTION_CA_FILE or TLS_OPTION_SERVER_NAME must be defined")

func StaticCertKeyContent() (cert, key []byte, err error) {
	certFile, keyFile := os.Getenv(StaticCertFile), os.Getenv(StaticKeyFile)
	if len(certFile) == 0 || len(keyFile) == 0 {
		return nil, nil, ErrCAFileNotInCluster
	}

	cert, err = os.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}
	key, err = os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func StaticServerCerts() (tls.Certificate, error) {
	certFile, keyFile := os.Getenv(StaticCertFile), os.Getenv(StaticKeyFile)
	if len(certFile) == 0 || len(keyFile) == 0 {
		return tls.Certificate{}, ErrCAFileNotInCluster
	}
	return tls.LoadX509KeyPair(certFile, keyFile)
}

func StaticClientCerts() (string, string, error) {
	caFile, serverName := os.Getenv(StaticCaFile), os.Getenv(StaticServerName)
	if len(caFile) == 0 || len(serverName) == 0 {
		return "", "", ErrGrpcTLSNotInCluster
	}
	return caFile, serverName, nil
}
