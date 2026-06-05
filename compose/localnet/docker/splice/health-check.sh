#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

wget --no-verbose --tries=1 --spider http://localhost:5012/api/scan/readyz
wget --no-verbose --tries=1 --spider http://localhost:5014/api/sv/readyz

# SV
wget --no-verbose --tries=1 --spider "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00/api/validator/readyz"
# Participant 1
wget --no-verbose --tries=1 --spider "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}01/api/validator/readyz"
# Participant 2
wget --no-verbose --tries=1 --spider "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}02/api/validator/readyz"
# Participant 3
wget --no-verbose --tries=1 --spider "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}03/api/validator/readyz"
# Participant 4
wget --no-verbose --tries=1 --spider "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}04/api/validator/readyz"
# Participant 5
wget --no-verbose --tries=1 --spider "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}05/api/validator/readyz"
