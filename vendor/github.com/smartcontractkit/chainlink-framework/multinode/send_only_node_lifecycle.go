package multinode

import (
	"fmt"
	"time"
)

// verifyLoop may only be triggered once, on Start, if initial chain ID check
// fails.
//
// It will continue checking until success and then exit permanently.
func (s *sendOnlyNode[CHAIN_ID, RPC]) verifyLoop() {
	defer s.wg.Done()
	ctx, cancel := s.chStop.NewCtx()
	defer cancel()

	backoff := NewRedialBackoff()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff.Duration()):
		}
		chainID, err := s.rpc.ChainID(ctx)
		if err != nil {
			ok := s.IfStarted(func() {
				if changed := s.setState(nodeStateUnreachable); changed {
					s.metrics.IncrementNodeTransitionsToUnreachable(ctx, s.name)
				}
			})
			if !ok {
				return
			}
			s.log.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
			continue
		} else if chainID.String() != s.chainID.String() {
			ok := s.IfStarted(func() {
				if changed := s.setState(nodeStateInvalidChainID); changed {
					s.metrics.IncrementNodeTransitionsToInvalidChainID(ctx, s.name)
				}
			})
			if !ok {
				return
			}
			s.log.Errorf(
				"sendonly rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
				chainID.String(),
				s.chainID.String(),
				s.name,
			)

			continue
		}
		ok := s.IfStarted(func() {
			if changed := s.setState(nodeStateAlive); changed {
				s.metrics.IncrementNodeTransitionsToAlive(ctx, s.name)
			}
		})
		if !ok {
			return
		}
		s.log.Infow("Sendonly RPC Node is online", "nodeState", s.state)
		return
	}
}
