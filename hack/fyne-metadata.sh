#!/bin/bash

set -e

script_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)"
base_dir="$(echo "${script_dir}" | xargs dirname)"

# shellcheck source=version.sh
source "${script_dir}/version.sh"

export gui_version="${RELEASE_VERSION#v}"

envsubst < "${base_dir}/templates/FyneApp.toml" > "${base_dir}/cmd/gui/FyneApp.toml"
