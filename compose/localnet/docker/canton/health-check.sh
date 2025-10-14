#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

# SV
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00" grpc.health.v1.Health/Check

# Participant 1
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}01"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}01" grpc.health.v1.Health/Check

# Participant 2
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}02"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}02" grpc.health.v1.Health/Check

# Participant 3
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}03"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}03" grpc.health.v1.Health/Check

# Participant 4
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}04"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}04" grpc.health.v1.Health/Check

# Participant 5
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}05"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}05" grpc.health.v1.Health/Check
