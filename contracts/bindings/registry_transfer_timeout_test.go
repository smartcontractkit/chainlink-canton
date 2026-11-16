package bindings_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	registryburnminttokenpool "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
)

// The factory decodes TransferTimeout with decodeTransferTimeoutAt in
// contracts/ccip/registry-pools/burn-mint-token-pool/daml/CCIP/Registry/BurnMintTokenPoolV2CodecGen.daml,
// which only accepts the single tag bytes 00 (Indefinite) and 01 (RelativeHours).
func TestDeployRegistryBurnMintTokenPoolParams_TransferTimeoutHex(t *testing.T) {
	hours := types.INT64(24)
	unit := types.UNIT{}

	tests := []struct {
		name    string
		timeout registryburnminttokenpool.TransferTimeout
		wantEnd string
	}{
		{
			name:    "RelativeHours",
			timeout: registryburnminttokenpool.TransferTimeout{RelativeHours: &hours},
			wantEnd: "010000000000000018",
		},
		{
			name:    "Indefinite",
			timeout: registryburnminttokenpool.TransferTimeout{Indefinite: &unit},
			wantEnd: "00",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := factory.DeployRegistryBurnMintTokenPoolParams{
				InstanceId:         "pool-1",
				InstrumentId:       splice_api_token_holding_v1.InstrumentId{Admin: "admin::1220", Id: "LINK"},
				Decimals:           10,
				Observers:          []types.PARTY{"observer::1220"},
				TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: "tar@admin::1220"},
				FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: "fq@admin::1220"},
				RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: "rmn@admin::1220"},
				PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
				TransferTimeout:    tc.timeout,
			}

			raw, err := params.MarshalHex()
			require.NoError(t, err)
			// MarshalHex returns raw bytes; the MCMS transport hex-encodes them.
			encoded := hex.EncodeToString([]byte(raw))
			require.Truef(t, strings.HasSuffix(encoded, tc.wantEnd),
				"encoding must end with tag %q, got tail %q", tc.wantEnd, encoded[max(0, len(encoded)-40):])
		})
	}
}
