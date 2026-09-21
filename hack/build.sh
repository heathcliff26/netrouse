#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

bin_dir="${base_dir}/bin"
script_dir="${base_dir}/hack"

GOOS="${GOOS:-$(go env GOOS)}"
GOARCH="${GOARCH:-$(go env GOARCH)}"

GO_LD_FLAGS="${GO_LD_FLAGS:-"-s"}"

# shellcheck source=version.sh
source "${script_dir}/version.sh"

echo "Building netrouse version ${RELEASE_VERSION}"
GO_LD_FLAGS+=" -X github.com/heathcliff26/netrouse/pkg/version.gitVersion=${RELEASE_VERSION}"

output_name="${bin_dir}/netrouse"
if [ "${1}" != "" ]; then
    output_name="${bin_dir}/${1}"
fi

if [ "${GOOS}" == "windows" ]; then
    output_name="${output_name}.exe"
fi

pushd "${base_dir}" >/dev/null

echo "Building $(basename "${output_name}")"
GOOS="${GOOS}" GOARCH="${GOARCH}" CGO_ENABLED=0 go build -ldflags="${GO_LD_FLAGS}" -o "${output_name}" ./cmd/cli/...
