//go:build !prodledger

package ledgertarget

import (
	pprof "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	receiverop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	senderop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"
)

var (
	CreateRouterOperation    = pprof.CreateRouter
	ReceiverDeployOperation  = receiverop.Deploy
	SenderDeployOperation    = senderop.Deploy
	ReceiverExecuteOperation = receiverop.Execute
	SenderSendOperation      = senderop.Send
)
