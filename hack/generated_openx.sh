#!/usr/bin/env bash

# set -o errexit
# set -o nounset
# set -o pipefail

export GOPATH=$(go env | grep -i gopath | awk '{split($0,a,"\""); print a[2]}')
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")

# The working directory which was the root path of our project.
ROOT_PACKAGE="github.com/kzz45/neverdown"
# API Group
CUSTOM_RESOURCE_NAME="$1"
# API Version
CUSTOM_RESOURCE_VERSION="$2"

deepcopy-gen \
    --go-header-file="${SCRIPT_ROOT}"/boilerplate.go.txt \
    --input-dirs "$ROOT_PACKAGE/pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION" \
    -O zz_generated.deepcopy \
    --bounding-dirs "$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION"

Packages="$ROOT_PACKAGE/pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION"

go-to-protobuf \
    --go-header-file="${SCRIPT_ROOT}"/boilerplate.go.txt \
    --apimachinery-packages="+k8s.io/apimachinery/pkg/util/intstr,+k8s.io/apimachinery/pkg/api/resource,+k8s.io/apimachinery/pkg/runtime/schema,+k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/apis/meta/v1beta1,k8s.io/apimachinery/pkg/apis/testapigroup/v1,-k8s.io/api/core/v1,-k8s.io/api/apps/v1,-k8s.io/api/autoscaling/v2" \
    --verify-only=false \
    --packages="${Packages}" \
    --proto-import="${GOPATH}/pkg/mod/" \
    --clean=false \
    --only-idl=false \
    --keep-gogoproto=false
