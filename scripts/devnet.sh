#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

IMAGE_NAME="ghcr.io/digital-asset/decentralized-canton-sync/docker/canton-participant"
IMAGE_VERSION="0.4.20"

DOCKER_ARGS=(
  --rm
  -it
  -p 127.0.0.1:1234:1234
  -p 127.0.0.1:2234:2234
  -p 127.0.0.1:5011:5011
  -p 127.0.0.1:5012:5012
  -p 127.0.0.1:5021:5021
  -p 127.0.0.1:5022:5022
  -v "$SCRIPT_DIR/devnet.bootstrap:/devnet.bootstrap"
  --entrypoint /app/bin/canton
  "${IMAGE_NAME}:${IMAGE_VERSION}"
  # sample config with 2 participants, the default from canton-open-source:
  # https://hub.docker.com/layers/digitalasset/canton-open-source/2.3.20/images/sha256-0966fbb8857aa60180850d3b7540c32a88d268684a94d247d16ea63237aa10bb
  --config examples/01-simple-topology/simple-topology.conf
  # expose admin and gRPC ledger APIs
  -C canton.participants.participant1.admin-api.address=0.0.0.0
  -C canton.participants.participant1.ledger-api.address=0.0.0.0
  -C canton.participants.participant2.admin-api.address=0.0.0.0
  -C canton.participants.participant2.ledger-api.address=0.0.0.0
  # expose and enable the JSON ledger API
  # https://github.com/digital-asset/canton/blob/92339b6f98faaecbe3adbfb71293ed9cbfb30204/community/ledger/ledger-json-api/src/main/scala/com/digitalasset/canton/http/HttpServerConfig.scala#L26
  -C canton.participants.participant1.http-ledger-api.server.address=0.0.0.0
  -C canton.participants.participant1.http-ledger-api.server.port=1234
  -C canton.participants.participant2.http-ledger-api.server.address=0.0.0.0
  -C canton.participants.participant2.http-ledger-api.server.port=2234
  --bootstrap /devnet.bootstrap
)

exec docker run "${DOCKER_ARGS[@]}"
