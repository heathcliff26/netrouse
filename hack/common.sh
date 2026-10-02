#!/bin/bash

# renovate: datasource=docker depName=ghcr.io/heathcliff26/go-fyne-ci extractVersion=^(?<version>.*)$
export BUILDER_VERSION=202609301715
export BUILDER_IMAGE="${BUILDER_IMAGE:-ghcr.io/heathcliff26/go-fyne-ci:$BUILDER_VERSION}"
