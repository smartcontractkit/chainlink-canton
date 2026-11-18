#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

# SV
echo "Checking ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}00"
wget -Y off -q --tries=1 --timeout=5 -O /dev/null "http://localhost:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}00/health/ready"

# Participant 1
echo "Checking ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}01"
wget -Y off -q --tries=1 --timeout=5 -O /dev/null "http://localhost:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}01/health/ready"

# Participant 2
echo "Checking ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}02"
wget -Y off -q --tries=1 --timeout=5 -O /dev/null "http://localhost:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}02/health/ready"

# Participant 3
echo "Checking ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}03"
wget -Y off -q --tries=1 --timeout=5 -O /dev/null "http://localhost:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}03/health/ready"

# Participant 4
echo "Checking ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}04"
wget -Y off -q --tries=1 --timeout=5 -O /dev/null "http://localhost:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}04/health/ready"

# Participant 5
echo "Checking ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}05"
wget -Y off -q --tries=1 --timeout=5 -O /dev/null "http://localhost:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}05/health/ready"
