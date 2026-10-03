#!/bin/sh

if [ "$1" != "0" ] && [ "$1" != "remove" ] && [ "$1" != "purge" ]; then
    exit 0
fi

echo "Clean up netroused service"
systemctl unmask netroused.service
systemctl stop netroused.service
systemctl disable netroused.service

