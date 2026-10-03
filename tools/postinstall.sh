#!/bin/sh

systemctl daemon-reload
systemctl enable --now netroused.service
