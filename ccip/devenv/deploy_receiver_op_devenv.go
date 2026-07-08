//go:build !prodledger

package devenv

import receiverop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"

var ccipReceiverDeployOperation = receiverop.Deploy
