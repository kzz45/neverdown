.PHONY: mod

WORK_PATH := $(shell pwd)
ROOT_PATH := $(shell echo ${HOME})
SECRET := "authx-secret"
PROJECT := "neverdown"
tag ?= "0.0.1"

# docker hub
DISCOVERY_APISERVER := "kongzz45/discovery-controlplane:$(tag)"
AUTHX_APISERVER := "kongzz45/authx-apiserver:$(tag)"
AUTHX_DASHBOARD := "kongzz45/authx-dashboard:$(tag)"
JINGX_APISERVER := "kongzz45/jingx-apiserver:$(tag)"
JINGX_DASHBOARD := "kongzz45/jingx-dashboard:$(tag)"
OPENX_APISERVER := "kongzz45/openx-apiserver:$(tag)"
OPENX_DASHBOARD := "kongzz45/openx-dashboard:$(tag)"

mod:
	go mod tidy

gen-clientset:
	cd hack && bash generate-internal-groups.sh \
	"client,conversion,deepcopy,defaulter,informer,lister" \
	github.com/kzz45/neverdown/pkg/client-go \
	github.com/kzz45/neverdown/pkg/apis \
	github.com/kzz45/neverdown/pkg/apis \
	"openx:v1"

gen-api-protobuf:
	cd hack && bash generated_openx.sh openx v1

gen-certs:
	cd certs && bash gen-certs.sh

gen-authx-frontend-protox:
	cd authx_frontend && ./node_modules/protobufjs-cli/bin/pbjs \
	-t static-module --es6 -w es6 -o src/proto/proto.js src/proto/protos/*.proto

run-authx-frontend-local:
	cd authx_frontend && npm run dev

build-authx-frontend:
	-i docker image rm $(AUTHX_DASHBOARD)
	cd authx_frontend && npm run build
	cp -r authx_frontend/dist .
	cp authx_frontend/Dockerfile . && docker build -t $(AUTHX_DASHBOARD) .
	rm -f Dockerfile && rm -rf dist

gen-jingx-frontend-protox:
	cd jingx_frontend && ./node_modules/protobufjs-cli/bin/pbjs \
	-t static-module --es6 -w es6 -o src/proto/proto.js src/proto/protos/*.proto

run-jingx-frontend-local:
	cd jingx_frontend && npm run dev

build-jingx-frontend:
	-i docker image rm $(JINGX_DASHBOARD)
	cd jingx_frontend && npm run build
	cp -r jingx_frontend/dist .
	cp jingx_frontend/Dockerfile.quick . && docker build -f Dockerfile.quick -t $(JINGX_DASHBOARD) .
	rm -f Dockerfile.quick && rm -rf dist

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
	go run cmd/openx/main.go -kubeconfig=$(ROOT_PATH)/.kube/config -listenPort=8080

build-openx:
	-i docker image rm $(OPENX_APISERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openx-apiserver cmd/openx/main.go
	cp cmd/openx/Dockerfile.quick . && docker build -f Dockerfile.quick -t $(OPENX_APISERVER) .
	rm -f Dockerfile.quick && rm -f openx-apiserver

gen-openx-frontend-protox:
	cd openx_frontend && ./node_modules/protobufjs/bin/pbjs \
	-t static-module --es6 -w es6 -o proto/proto.js proto/protos/*.proto

run-openx-frontend-local:
	cd openx_frontend && npm run dev

build-openx-frontend:
	-i docker image rm $(OPENX_DASHBOARD)
	cd openx_frontend && npm run build
	cp -r openx_frontend/dist .
	cp openx_frontend/Dockerfile . && docker build -t $(OPENX_DASHBOARD) .
	rm -f Dockerfile && rm -rf dist
