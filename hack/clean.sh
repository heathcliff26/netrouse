#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

folders=("bin" "dist" "coverprofiles" ".idea")
files=("coverprofile.out" "cmd/gui/FyneApp.toml" "cmd/gui/fyne_metadata_init.go" "netrouse.iml" "debug.keystore")

for folder in "${folders[@]}"; do
    if ! [ -e "${base_dir}/${folder}" ]; then
        continue
    fi
    echo "Removing ${folder}"
    rm -rf "${base_dir:-.}/${folder}"
done

for file in "${files[@]}"; do
    if ! [ -e "${base_dir}/${file}" ]; then
        continue
    fi
    echo "Removing ${file}"
    rm "${base_dir:-.}/${file}"
done
