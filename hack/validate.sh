#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

echo "Check if source code is formatted"
make fmt
rc=0
git update-index --refresh && git diff-index --quiet HEAD -- || rc=1
if [ $rc -ne 0 ]; then
    echo "FATAL: Need to run \"make fmt\""
    exit 1
fi

echo "Check if the bootstrap file is generated"
make generate-bootstrap
rc=0
git update-index --refresh && git diff-index --quiet HEAD -- || rc=1
if [ $rc -ne 0 ]; then
    echo "FATAL: Need to run \"make generate-bootstrap\""
    exit 1
fi

echo "Check if the swagger docs are generated"
make generate-swagger
rc=0
git update-index --refresh && git diff-index --quiet HEAD -- || rc=1
if [ $rc -ne 0 ]; then
    echo "FATAL: Need to run \"make generate-swagger\""
    exit 1
fi

echo "All files are up to date"

if grep "/bootstrap/" "${base_dir}/static/index.html"; then
    echo "FATAL: Using files from the bootstrap folder inside index.html"
    exit 1
fi

echo "Check if metainfo file is valid"
make validate-metainfo
