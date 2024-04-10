.PHONY: mod

CUR_PATH := $(shell pwd)
ROOT_PATH := $(shell echo ${HOME})
SECRET := "authx-secret"
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

gen-authx-grpc-proto:
	protoc --go_out=pkg/authx/grpc --go-grpc_out=pkg/authx/grpc pkg/authx/grpc/proto/*.proto

gen-authx-http-proto:
	cd hack && bash authx_proto.sh

gen-authx-frontend-protox:
	cd authx_frontend && ./node_modules/protobufjs-cli/bin/pbjs \
	-t static-module --es6 -w es6 -o src/proto/proto.js src/proto/protos/*.proto

run-authx-local:
	DISCOVERY_SERVICE_HOST="127.0.0.1" \
	DISCOVERY_SERVICE_PORT=9443 \
	DISCOVERY_SERVICE_CAFILE=${CUR_PATH}/certs/ca.crt \
	TLS_OPTION_CERT_FILE=${CUR_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${CUR_PATH}/certs/server.key \
	AUTHX_SECRET="$(SECRET)" \
	TOKEN_EXPIRATION=36000 \
	go run ./cmd/authx/main.go

run-authx-frontend-local:
	cd authx_frontend && npm run dev

gen-jingx-http-proto:
	cd hack && bash jingx_proto.sh

run-jingx-local:
	DISCOVERY_SERVICE_HOST="127.0.0.1" \
	DISCOVERY_SERVICE_PORT=9443 \
	DISCOVERY_SERVICE_CAFILE=${CUR_PATH}/certs/ca.crt \
	AUTHORITY_SERVICE_HOST="127.0.0.1" \
	AUTHORITY_SERVICE_PORT=9443 \
	AUTHORITY_SERVICE_CAFILE=${CUR_PATH}/certs/ca.crt \
	TLS_OPTION_CERT_FILE=${CUR_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${CUR_PATH}/certs/server.key \
	go run ./cmd/jingx/main.go -listenPort=8083
