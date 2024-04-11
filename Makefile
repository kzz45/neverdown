.PHONY: mod

WORK_PATH := $(shell pwd)
ROOT_PATH := $(shell echo ${HOME})
SECRET := "authx-secret"
PROJECT := "neverdown"
TAG ?= "0.0.1"

# docker hub
DISCOVERY_APISERVER := "kongzz45/discovery:${TAG}"
AUTHX_APISERVER := "kongzz45/authx-apiserver:${TAG}"
JINGX_APISERVER := "kongzz45/jingx-apiserver:${TAG}"
OPENX_APISERVER := "kongzz45/openx-apiserver:${TAG}"

mod:
	go mod tidy

gen-clientset:
	cd hack && bash generate-internal-groups.sh \
	"client,conversion,deepcopy,defaulter,informer,lister" \
	github.com/kzz45/neverdown/pkg/client-go \
	github.com/kzz45/neverdown/pkg/apis \
	github.com/kzz45/neverdown/pkg/apis \
	"rbac:v1 audit:v1 jingx:v1 openx:v1"

gen-api-protobuf:
	cd hack && bash generated.sh rbac v1
	cd hack && bash generated.sh audit v1 "-github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	cd hack && bash generated.sh jingx v1 "-github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	cd hack && bash generated_openx.sh openx v1

gen-certs:
	cd certs && bash gen-certs.sh

run-discovery-local:
	TLS_OPTION_CERT_FILE=${WORK_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${WORK_PATH}/certs/server.key \
	ETCD_PREFIX="/registry" \
	go run ./cmd/discovery/main.go -listenPort=9443

build-discovery:
	docker image rm $(DISCOVERY_APISERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o discovery-controlplane cmd/discovery/main.go
	cp cmd/discovery/Dockerfile . && docker build -t $(DISCOVERY_APISERVER) .
	rm -f Dockerfile && rm -f discovery-controlplane

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
	DISCOVERY_SERVICE_CAFILE=${WORK_PATH}/certs/ca.crt \
	TLS_OPTION_CERT_FILE=${WORK_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${WORK_PATH}/certs/server.key \
	AUTHX_SECRET="$(SECRET)" \
	TOKEN_EXPIRATION=36000 \
	go run ./cmd/authx/main.go

run-authx-frontend-local:
	cd authx_frontend && npm run dev

gen-jingx-http-proto:
	cd hack && bash jingx_proto.sh

gen-jingx-frontend-protox:
	cd jingx_frontend && ./node_modules/protobufjs-cli/bin/pbjs \
	-t static-module --es6 -w es6 -o src/proto/proto.js src/proto/protos/*.proto

run-jingx-local:
	DISCOVERY_SERVICE_HOST="127.0.0.1" \
	DISCOVERY_SERVICE_PORT=9443 \
	DISCOVERY_SERVICE_CAFILE=${WORK_PATH}/certs/ca.crt \
	AUTHORITY_SERVICE_HOST="127.0.0.1" \
	AUTHORITY_SERVICE_PORT=9443 \
	AUTHORITY_SERVICE_CAFILE=${WORK_PATH}/certs/ca.crt \
	TLS_OPTION_CERT_FILE=${WORK_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${WORK_PATH}/certs/server.key \
	go run ./cmd/jingx/main.go -listenPort=8083

run-jingx-frontend-local:
	cd jingx_frontend && npm run dev

gen-openx-proto:
	cd hack && bash openx_proto.sh

run-openx-local:
	DISCOVERY_SERVICE_HOST="127.0.0.1" \
	DISCOVERY_SERVICE_PORT=9443 \
	DISCOVERY_SERVICE_CAFILE=${WORK_PATH}/certs/ca.crt  \
	AUTHORITY_SERVICE_HOST="127.0.0.1" \
	AUTHORITY_SERVICE_PORT=9443 \
	TOKEN_EXPIRATION=36000 \
	AUTHORITY_SERVICE_CAFILE=${WORK_PATH}/certs/ca.crt \
	TLS_OPTION_CERT_FILE=${WORK_PATH}/certs/server.crt \
	TLS_OPTION_KEY_FILE=${WORK_PATH}/certs/server.key \
	go run cmd/openx/main.go -kubeconfig=$(WORK_PATH)/config/kind.config -listenPort=8080

gen-openx-frontend-protox:
	cd openx_frontend && ./node_modules/protobufjs/bin/pbjs \
	-t static-module --es6 -w es6 -o proto/proto.js proto/protos/*.proto

run-openx-frontend-local:
	cd openx_frontend && npm run dev

OPENX_DASHBOARD_IMAGE := "kongzz45/openx-dashboard:$(TAG)"

build-openx-frontend:
	cd openx_frontend
	npm run build
	rm -rf dist/config
	docker build -t $(OPENX_DASHBOARD_IMAGE) .
	rm -rf dist
