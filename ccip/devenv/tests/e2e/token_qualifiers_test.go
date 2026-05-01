package canton

import (
	"testing"

	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
)

func burnMint20ToLockRelease20TokenQualifier(t *testing.T) string {
	t.Helper()

	for _, combo := range devenvcommon.AllTokenCombinations() {
		local := combo.LocalPoolAddressRef()
		remote := combo.RemotePoolAddressRef()
		if string(local.Type) == devenvcommon.BurnMintTokenPoolType &&
			local.Version.String() == "2.0.0" &&
			string(remote.Type) == devenvcommon.LockReleaseTokenPoolType &&
			remote.Version.String() == "2.0.0" {
			return local.Qualifier
		}
	}

	t.Fatal("could not find BurnMint 2.0.0 to LockRelease 2.0.0 token qualifier")
	return ""
}
