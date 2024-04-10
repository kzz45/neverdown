.PHONY: dev

# HARBOR_DOMAIN := xxx
PROJECT := neverdown
tag ?= 0.0.1
# IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/jingx-dashboard:$(tag)"
IMAGE := "$(PROJECT)/jingx-dashboard:$(tag)"

dev:
	npm run dev

web:
	docker buildx build -f Dockerfile -t $(IMAGE) .
	# docker push $(IMAGE)

web-quick:
	npm run build
	docker buildx build -f Dockerfile.quick -t $(IMAGE) .
	# docker push $(IMAGE)
	rm -rf dist

protox:
	./node_modules/protobufjs-cli/bin/pbjs \
	-t static-module --es6 -w es6 -o src/proto/proto.js src/proto/protos/*.proto
