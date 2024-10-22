#!/usr/bin/env bash

ROOT_PATH="$(cd "$(dirname "$0")/.." && pwd -P)"
BUILD_PATH="${ROOT_PATH}/build"

source "${BUILD_PATH}/init.sh"

cleanup() {
  rm -r ${BUILD_OUTPUT_PATH}
}

cleanup
