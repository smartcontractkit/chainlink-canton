#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

generate_jwt() {
  local sub="$1"
  local aud="$2"
  jwt-cli encode hs256 --s unsafe --p '{"sub": "'"$sub"'", "aud": "'"$aud"'"}'
}

SV_USER_TOKEN=$(generate_jwt "$API_USER_NAME" "$API_AUDIENCE")
export SV_USER_TOKEN

PARTICIPANT1_USER_TOKEN=$(generate_jwt "$API_USER_NAME" "$API_AUDIENCE")
export PARTICIPANT1_USER_TOKEN

PARTICIPANT2_USER_TOKEN=$(generate_jwt "$API_USER_NAME" "$API_AUDIENCE")
export PARTICIPANT2_USER_TOKEN

PARTICIPANT3_USER_TOKEN=$(generate_jwt "$API_USER_NAME" "$API_AUDIENCE")
export PARTICIPANT3_USER_TOKEN

PARTICIPANT4_USER_TOKEN=$(generate_jwt "$API_USER_NAME" "$API_AUDIENCE")
export PARTICIPANT4_USER_TOKEN

PARTICIPANT5_USER_TOKEN=$(generate_jwt "$API_USER_NAME" "$API_AUDIENCE")
export PARTICIPANT5_USER_TOKEN

# source all scripts from /app/pre-startup/on so that env variables exported by them are available in the current shell
for script in /app/pre-startup/on/*.sh; do
# shellcheck disable=SC1090
  [ -f "$script" ] && source "$script"
done
/app/bin/canton --no-tty -c /app/app.conf
