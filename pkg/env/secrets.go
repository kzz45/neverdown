package env

import (
	"errors"
	"os"
)

const (
	StaticAuthSecret = "AUTHX_SECRET"
)

var ErrGrpcAuthSecretNotInCluster = errors.New("unable to load grpc authentication configuration, AUTHX_SECRET must be defined")

func GrpcAuthenticationSecret() (string, error) {
	sc := os.Getenv(StaticAuthSecret)
	if len(sc) == 0 {
		return "", ErrGrpcAuthSecretNotInCluster
	}
	return sc, nil
}
