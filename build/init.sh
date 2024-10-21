#!/usr/bin/env bash

BUILD_OUTPUT_PATH="${ROOT_PATH}/_output"
BUILD_CMD_PATH="${BUILD_OUTPUT_PATH}/bin"

PACKAGE_NAME="github.com/bamboo-firewall/be"

VERSION="$(git describe --abbrev=0 --tags)"
BRANCH="$(git rev-parse --abbrev-ref HEAD)"
BUILD_TIME="$(date +%Y-%m-%dT%H:%M:%S%z)"

ORGANIZATION="ATAOCloud"

LDFLAGS="-s -w -X ${PACKAGE_NAME}/buildinfo.Version=${VERSION} \
		 -X ${PACKAGE_NAME}/buildinfo.GitBranch=${BRANCH} \
		 -X ${PACKAGE_NAME}/buildinfo.BuildDate=${BUILD_TIME} \
		 -X ${PACKAGE_NAME}/buildinfo.Organization=${ORGANIZATION}"

CGO_ENABLED=0

GO_BUILD="go build -buildvcs=false -a -installsuffix cgo"

source "${BUILD_PATH}/golang.sh"
