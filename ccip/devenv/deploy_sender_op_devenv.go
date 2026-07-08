//go:build !prodledger

package devenv

import senderop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"

var ccipSenderDeployOperation = senderop.Deploy
