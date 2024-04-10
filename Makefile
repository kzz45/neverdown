.PHONY: mod

CUR_PATH := $(shell pwd)
ROOT_PATH := $(shell echo ${HOME})
SECRET := "demo"
PROJECT := "neverdown"
TAG := "latest"

DISCOVERY_APISERVER := "${PROJECT}/discovery:${TAG}"
AUTHX_APISERVER := "${PROJECT}/authx-apiserver:${TAG}"
JINGX_APISERVER := "${PROJECT}/jingx-apiserver:${TAG}"
OPENX_APISERVER := "${PROJECT}/openx-apiserver:${TAG}"

mod:
	go mod tidy

gen-clientset:
	cd hack && bash generate-internal-groups.sh \
	"client,conversion,deepcopy,defaulter,informer,lister" \
	github.com/kzz45/neverdown/pkg/client-go \
	github.com/kzz45/neverdown/pkg/apis \
	github.com/kzz45/neverdown/pkg/apis \
	"rbac:v1 audit:v1 jingx:v1"

gen-api-protobuf:
	cd hack && bash generated.sh rbac v1
	cd hack && bash generated.sh audit v1 "-github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	cd hack && bash generated.sh jingx v1 "-github.com/kzz45/neverdown/pkg/apis/rbac/v1"

gen-certs:
	cd certs && bash gen-certs.sh

run-discovery-local:
	TLS_OPTION_CERT_FILE=${CUR_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${CUR_PATH}/certs/server.key \
	ETCD_PREFIX="/registry" \
	go run ./cmd/discovery/main.go -listenPort=9443
