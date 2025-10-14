#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

curl -f http://localhost:5012/api/scan/readyz
curl -f http://localhost:5014/api/sv/readyz

# SV
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00/api/validator/readyz"
# Participant 1
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}01/api/validator/readyz"
# Participant 2
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}02/api/validator/readyz"
# Participant 3
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}03/api/validator/readyz"
# Participant 4
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}04/api/validator/readyz"
# Participant 5
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}05/api/validator/readyz"
