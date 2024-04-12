.PHONY: mod

WORK_PATH := $(shell pwd)
ROOT_PATH := $(shell echo ${HOME})
SECRET := "authx-secret"
PROJECT := "neverdown"
tag ?= "0.0.1"

# docker hub
DISCOVERY_APISERVER := "kongzz45/discovery-controlplane:$(tag)"
AUTHX_APISERVER := "kongzz45/authx-apiserver:$(tag)"
AUTHX_DASHBOARD_IMAGE := "kongzz45/authx-dashboard:$(tag)"
JINGX_APISERVER := "kongzz45/jingx-apiserver:$(tag)"
JINGX_DASHBOARD_IMAGE := "kongzz45/jingx-dashboard:$(tag)"
OPENX_APISERVER := "kongzz45/openx-apiserver:$(tag)"
OPENX_DASHBOARD_IMAGE := "kongzz45/openx-dashboard:$(tag)"

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
	-i docker image rm $(DISCOVERY_APISERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o discovery-controlplane cmd/discovery/main.go
	cp cmd/discovery/Dockerfile.quick . && docker build -f Dockerfile.quick -t $(DISCOVERY_APISERVER) .
	rm -f Dockerfile.quick && rm -f discovery-controlplane

run-discovery-docker:
	docker run -d -p 9443:9443 -v ${PWD}/certs:/certs \
	-e TLS_OPTION_CERT_FILE=/certs/server.crt \
	-e TLS_OPTION_KEY_FILE=/certs/server.key \
	-e ETCD_ENDPOINTS="http://$(shell docker inspect --format '{{ .NetworkSettings.IPAddress }}' etcd-server):2379" \
	--name=discovery-controlplane kongzz45/discovery-controlplane:0.0.1
	

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

run-authx-docker:
	docker run -d -p 8087:8087 -v ${PWD}/certs:/certs \
	-e TLS_OPTION_CERT_FILE=/certs/server.crt \
	-e TLS_OPTION_KEY_FILE=/certs/server.key \
	-e DISCOVERY_SERVICE_HOST="$(shell docker inspect --format '{{ .NetworkSettings.IPAddress }}' discovery-controlplane)" \
	-e DISCOVERY_SERVICE_PORT=9443 \
	-e DISCOVERY_SERVICE_CAFILE=/certs/ca.crt \
	-e AUTHX_SECRET="$(SECRET)" \
	-e TOKEN_EXPIRATION=36000 \
	--name=authx-apiserver kongzz45/authx-apiserver:0.0.1

build-authx:
	-i docker image rm $(AUTHX_APISERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authx-apiserver cmd/authx/main.go
	cp cmd/authx/Dockerfile.quick . && docker build -f Dockerfile.quick -t $(AUTHX_APISERVER) .
	rm -f Dockerfile.quick && rm -f authx-apiserver

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

build-jingx:
	-i docker image rm $(JINGX_APISERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jingx-apiserver cmd/jingx/main.go
	cp cmd/jingx/Dockerfile.quick . && docker build -f Dockerfile.quick -t $(JINGX_APISERVER) .
	rm -f Dockerfile.quick && rm -f jingx-apiserver

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

build-openx:
	-i docker image rm $(OPENX_APISERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openx-apiserver cmd/openx/main.go
	cp cmd/openx/Dockerfile . && docker build -t $(OPENX_APISERVER) .
	rm -f Dockerfile && rm -f openx-apiserver

gen-openx-frontend-protox:
	cd openx_frontend && ./node_modules/protobufjs/bin/pbjs \
	-t static-module --es6 -w es6 -o proto/proto.js proto/protos/*.proto

run-openx-frontend-local:
	cd openx_frontend && npm run dev

build-openx-frontend:
	cd openx_frontend
	npm run build
	rm -rf dist/config
	docker build -t $(OPENX_DASHBOARD_IMAGE) .
	rm -rf dist
