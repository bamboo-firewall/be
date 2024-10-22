#!/usr/bin/env bash

readonly SUPPORTED_PLATFORMS=(
  linux/amd64
  linux/arm64
  darwin/amd64
  darwin/arm64
  windows/amd64
  windows/arm64
)

golang::build_binaries() {
  local -a platforms
  if [[ "$3" == "all" ]]; then
    platforms=("${SUPPORTED_PLATFORMS[@]}")
  else
    local host_platform
    host_platform=$(golang::host_platform)
    platforms+=("${host_platform}")
  fi

  for platform in "${platforms[@]}"; do
    golang::build_binary_for_platform ${platform} $1 $2
  done
}

golang::build_binary_for_platform() {
  local platform="$1"
  local build_dir="$2"
  local bin_name="$3"

  GOOS=${platform%%/*}
  GOARCH=${platform##*/}
  output="${BUILD_CMD_PATH}/${GOOS}/${GOARCH}/${bin_name}"

  CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} ${GO_BUILD} -ldflags="${LDFLAGS}" -o ${output} ${build_dir}
}

golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}
