.PHONY: dev

DOMAIN := "dockerhub"
PROJECT := "neverdown"
TAG ?= 0.0.1
# IMAGE := "${DOMAIN}/${PROJECT}/authx-dashboard:${TAG}"
IMAGE := "${PROJECT}/authx-dashboard:${TAG}"

dev:
	npm run dev

web:
	npm run build:prod
	docker build -t ${IMAGE} .
	bash hack/jingx.sh cli ${PROJECT} authx-dashboard ${TAG}

protox:
	./node_modules/protobufjs-cli/bin/pbjs \
	-t static-module --es6 -w es6 -o src/proto/proto.js src/proto/protos/*.proto

clean:
	$(shell docker images|grep ${IMAGE}|awk '{print $3}'|xargs docker rmi -f)
