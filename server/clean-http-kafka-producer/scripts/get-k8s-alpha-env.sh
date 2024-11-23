#!/usr/bin/env bash

COMMON_SECRET_NAME="domain-common-secret"

echo "$(kubectl get secret ${COMMON_SECRET_NAME} -n alpha -o json | jq 'select(.data != null) |.data| del(.SENTRY_DSN) | to_entries | map("\(.key)=\(.value | @base64d)") | .[]' --raw-output)"
