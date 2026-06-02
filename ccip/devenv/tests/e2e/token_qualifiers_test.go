package canton

import (
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

func defaultDevenvTokenLane(t *testing.T, lib ccv.Lib, in *ccv.Cfg, srcSelector, destSelector uint64) devenvtests.TokenLane {
	t.Helper()
	return devenvtests.DefaultDevenvTokenLane(t, in, lib, srcSelector, destSelector)
}
