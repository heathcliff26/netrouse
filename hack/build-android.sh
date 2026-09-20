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

if [ -z "${KEYSTORE}" ]; then
    echo "No keystore specified, using debug keystore"
    export KEYSTORE="debug.keystore"
    export KEYSTORE_PASS="android"
    export KEYSTORE_ALIAS="androiddebugkey"
    if [ ! -e "${KEYSTORE}" ]; then
        keytool -genkeypair -v \
            -keystore "${KEYSTORE}" \
            -storepass "${KEYSTORE_PASS}" \
            -alias "${KEYSTORE_ALIAS}" \
            -keypass android \
            -keyalg RSA \
            -keysize 2048 \
            -validity 10000 \
            -dname "CN=NetRouse Debug,O=NetRouse,C=DE"
    fi
fi

bin/fyne release --os android --app-build 1 \
    --keystore "${KEYSTORE}" \
    --keystore-pass "${KEYSTORE_PASS}" \
    --key-name "${KEYSTORE_ALIAS}"

bin/fyne package --os android --release --app-build 1
apksigner sign --ks "${KEYSTORE}" --ks-pass "pass:${KEYSTORE_PASS}" --ks-key-alias "${KEYSTORE_ALIAS}" \
    --v1-signing-enabled true \
    --v2-signing-enabled true \
    --v3-signing-enabled true \
    NetRouse.apk

mkdir -p dist
mv NetRouse.aab NetRouse.apk dist/
rm NetRouse.apk.idsig

for file in "${files[@]}"; do
    rm "${file}"
done
