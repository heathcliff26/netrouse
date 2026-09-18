#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

git_version="$(git describe --tags --always --dirty)"

export RELEASE_VERSION="${RELEASE_VERSION:-$git_version}"

envsubst < "${base_dir}/templates/FyneApp.toml" > "${base_dir}/cmd/gui/FyneApp.toml"
