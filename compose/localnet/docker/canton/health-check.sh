#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

# SV
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00"
grpc-health-probe -addr="localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00"

# Participant 1
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}01"
grpc-health-probe -addr="localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}01"

# Participant 2
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}02"
grpc-health-probe -addr="localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}02"

# Participant 3
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}03"
grpc-health-probe -addr="localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}03"

# Participant 4
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}04"
grpc-health-probe -addr="localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}04"

# Participant 5
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}05"
grpc-health-probe -addr="localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}05"
