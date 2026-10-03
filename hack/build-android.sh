#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

cd "${base_dir}"

# shellcheck source=version.sh
source hack/version.sh

if [[ ! "${RELEASE_VERSION}" =~ ^v?([0-9]+)\.([0-9]+)\.([0-9]+)([-+].*)?$ ]]; then
    echo "Invalid release version: ${RELEASE_VERSION}" >&2
    exit 1
fi

version_major=$((10#${BASH_REMATCH[1]}))
version_minor=$((10#${BASH_REMATCH[2]}))
version_patch=$((10#${BASH_REMATCH[3]}))
if (( version_major >= 2100 || version_minor >= 1000 || version_patch >= 1000 )); then
    echo "Release version is too large for an Android version code: ${RELEASE_VERSION}" >&2
    exit 1
fi
app_build=$((version_major * 1000000 + version_minor * 1000 + version_patch))
if (( app_build < 1 )); then
    echo "Android version code must be positive: ${RELEASE_VERSION}" >&2
    exit 1
fi

hack/fyne-metadata.sh

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
            -keyalg RSA \
            -keysize 2048 \
            -validity 10000 \
            -dname "CN=NetRouse Debug,O=NetRouse,C=DE"
    fi
fi

bin/fyne release --os android --app-build "${app_build}" \
    --keystore "${KEYSTORE}" \
    --keystore-pass "${KEYSTORE_PASS}" \
    --key-name "${KEYSTORE_ALIAS}"

mkdir -p dist
mv NetRouse.aab dist/

targets=("universal" "arm" "arm64" "amd64")

for target in "${targets[@]}"; do
    os="android"
    if [ "${target}" != "universal" ]; then
        os="android/${target}"
    fi
    echo "Building for ${target}"
    bin/fyne package --os "${os}" --release --app-build "${app_build}"
    apksigner sign \
        --ks "${KEYSTORE}" \
        --ks-pass "pass:${KEYSTORE_PASS}" \
        --ks-key-alias "${KEYSTORE_ALIAS}" \
        --v1-signing-enabled true \
        --v2-signing-enabled true \
        --v3-signing-enabled true \
        NetRouse.apk
    mv NetRouse.apk dist/NetRouse-"${target}".apk
done

rm NetRouse.apk.idsig
