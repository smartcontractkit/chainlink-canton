package factory

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	registryratelimiter "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
)

// encodeDamlText mirrors MCMS.Codec.encodeText and encodeParty: a 1-byte length prefix, then UTF-8 bytes.
func encodeDamlText(s string) string {
	return fmt.Sprintf("%02x", len(s)) + hex.EncodeToString([]byte(s))
}

func TestEncodeDeployRegistryRateLimiter_matchesFactoryDecoderLayout(t *testing.T) {
	t.Parallel()

	encoded, err := encodeDeployRegistryRateLimiter(factorybindings.DeployRegistryRateLimiter{
		Contract: registryratelimiter.RateLimiter{
			InstanceId:          "rl-1",
			PoolInstanceId:      "pool-1",
			PoolOwner:           "owner::1220ab",
			RemoteChainSelector: "137",
			Direction:           registryratelimiter.RateLimitDirectionRateLimitDirection_Inbound,
			Mode:                registryratelimiter.RateLimitModeRateLimitMode_CustomFinality,
			IsEnabled:           true,
			Capacity:            "1000",
			Rate:                "10",
			Tokens:              "1000",
			Observers:           []types.PARTY{"eds::1220cd"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "DeployRegistryRateLimiter", encoded.Choice)

	// Field order and tags follow decodeDeployRegistryRateLimiterParams in FactoryV2.daml.
	want := encodeDamlText("rl-1") +
		encodeDamlText("pool-1") +
		encodeDamlText("137") +
		"01" + // RateLimitDirection_Inbound
		"01" + // RateLimitMode_CustomFinality
		"01" + // isEnabled
		encodeDamlText("1000") +
		encodeDamlText("10") +
		"01" + encodeDamlText("eds::1220cd")
	require.Equal(t, want, encoded.OperationData)
}

func TestEncodeDeployRegistryBurnMintTokenPool_encodesTransferTimeoutTag(t *testing.T) {
	t.Parallel()

	owner := types.PARTY("owner::1220ab")
	encoded, err := encodeDeployRegistryBurnMintTokenPool(factorybindings.DeployRegistryBurnMintTokenPool{
		Contract: registrybnm.BurnMintTokenPool{
			InstanceId:         "registry-bmtp-1",
			PoolOwner:          owner,
			CcipOwner:          owner,
			InstrumentId:       splice_api_token_holding_v1.InstrumentId{Admin: owner, Id: "LINK"},
			Decimals:           10,
			Observers:          []types.PARTY{owner},
			PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
			TransferTimeout:    registrybnm.TransferTimeout{RelativeHours: new(types.INT64(24))},
			Deps: registrybnm.BurnMintTokenPoolDeps{
				TokenAdminRegistry: chainlinkapi.RawInstanceAddress{Unpack: "tar@owner::1220ab"},
				RmnRemote:          chainlinkapi.RawInstanceAddress{Unpack: "rmn@owner::1220ab"},
				FeeQuoter:          chainlinkapi.RawInstanceAddress{Unpack: "fq@owner::1220ab"},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "DeployRegistryBurnMintTokenPool", encoded.Choice)

	// transferTimeout is the last field: tag 1 (RelativeHours) then int64 24, per decodeTransferTimeoutAt.
	require.True(t, strings.HasSuffix(encoded.OperationData, "01"+"0000000000000018"), encoded.OperationData)
}
