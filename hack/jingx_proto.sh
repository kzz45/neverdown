#!/usr/bin/env bash

export GOPATH=$(go env | grep -i gopath | awk '{split($0,a,"\""); print a[2]}')
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")

GENS="$1"

# The working directory which was the root path of our project.
ROOT_PACKAGE="github.com/kzz45/neverdown"
# API Group
CUSTOM_RESOURCE_NAME="jingx"
# API Version
CUSTOM_RESOURCE_VERSION="proto"

Packages="$ROOT_PACKAGE/pkg/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION"

# protobuf
go-to-protobuf \
    --go-header-file="${GOPATH}/src/github.com/kzz45/neverdown/hack/boilerplate.go.txt" \
    --apimachinery-packages="-k8s.io/apimachinery/pkg/apis/meta/v1,-k8s.io/apimachinery/pkg/runtime/schema,-github.com/kzz45/neverdown/pkg/apis/rbac/v1" \
    --verify-only=false \
    --packages="${Packages}" \
    --proto-import="${GOPATH}/pkg/mod/" \
    --clean=false \
    --only-idl=false \
    --keep-gogoproto=false
