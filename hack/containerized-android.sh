#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

# shellcheck source=common.sh
source "${base_dir}/hack/common.sh"

export BUILDER_IMAGE_ANDROID="${BUILDER_IMAGE_ANDROID:-${BUILDER_IMAGE}-android}"
export BUILDER_IMAGE="${BUILDER_IMAGE_ANDROID}"

"${base_dir}/hack/containerized" hack/build-android.sh
