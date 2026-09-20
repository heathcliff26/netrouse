#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

cd "${base_dir}"

hack/fyne-metadata.sh

files=("main.go" "log.go" "FyneApp.toml")

for file in "${files[@]}"; do
    cp "cmd/gui/${file}" "${file}"
done

sed -i 's#../../##g' FyneApp.toml

bin/fyne package --os android --release

mkdir -p dist
mv ./*.apk dist/

for file in "${files[@]}"; do
    rm "${file}"
done
