#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"
dist_dir="${base_dir}/dist"
name="$(yq -r '.project_name' "${base_dir}/.goreleaser.yaml")"

echo "Preparing metadata for fyne"
"${base_dir}/hack/fyne-metadata.sh"

echo "Building releaser artifacts with goreleaser"
goreleaser release --skip=announce,publish,validate --clean -p 1

echo "Moving release artifacts to top level of dist directory"
artifacts="$(cat "${dist_dir}/artifacts.json" | jq -r -c '.[]')"
echo "${artifacts}" | while read -r artifact; do
    if ! [[ "$(echo "${artifact}" | jq -r '.name')" =~ ^${name}(\.exe)?$ ]]; then
        continue
    fi

    path="${base_dir}/$(echo "${artifact}" | jq -r '.path')"
    goarch="$(echo "${artifact}" | jq -r '.goarch')"
    ext="$(echo "${artifact}" | jq -r '.extra.Ext')"

    mv "${path}" "${dist_dir}/${name}-${goarch}${ext}"

    path="$(dirname "${path}")"
    rm -r "${path}"
done


echo "Cleaning up dist directory"
rm -r "${dist_dir}/artifacts.json" "${dist_dir}/config.yaml" "${dist_dir}/metadata.json" "${dist_dir}"/netrouse_*_checksums.txt "${dist_dir}"/gui_linux_*
