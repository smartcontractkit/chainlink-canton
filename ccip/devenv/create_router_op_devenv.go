//go:build !prodledger

package devenv

import pprof "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"

var createRouterOperation = pprof.CreateRouter
