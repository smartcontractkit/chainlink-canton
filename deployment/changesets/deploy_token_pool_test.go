package changesets

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/auth"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestDeployTokenPool deploys TAR, then runs DeployTokenPool (deploy pool + register with TAR in one changeset),
// then verifies TAR config by querying GetTokenConfig and asserting tokenPoolOwner is set.
func TestDeployTokenPool(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	token, err := bc.(*canton.Chain).Participants[0].JWTProvider.Token(t.Context())
	require.NoError(t, err)

	insecureCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	adminApiClient, err := grpc.NewClient(bc.(*canton.Chain).Participants[0].Endpoints.AdminAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	require.NoError(t, err, "Failed to dial gRPC admin API")
	ledgerApiClient, err := grpc.NewClient(bc.(*canton.Chain).Participants[0].Endpoints.GRPCLedgerAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	require.NoError(t, err, "Failed to dial gRPC ledger API")

	packageServiceClient := participantv30.NewPackageServiceClient(adminApiClient)
	userManagementServiceClient := admin.NewUserManagementServiceClient(ledgerApiClient)

	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	tarDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	poolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	_, err = packageServiceClient.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: commonDar},
			{Bytes: tarDar},
			{Bytes: poolDar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload dar files")

	userResp, err := userManagementServiceClient.GetUser(t.Context(), &admin.GetUserRequest{
		UserId: "user-participant1",
	})
	require.NoError(t, err, "failed to get user")
	user := userResp.GetUser()
	party := user.GetPrimaryParty()

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)
	cantonChain := bc.(*canton.Chain)
	deps := dependencies.CantonDeps{
		Chain:                *cantonChain,
		CommandServiceClient: apiv2.NewCommandServiceClient(ledgerApiClient),
		StateServiceClient:   apiv2.NewStateServiceClient(ledgerApiClient),
		Party:                party,
	}

	// Deploy TAR so we have an instance address for register-with-TAR
	tarAddrRef, err := cld_ops.ExecuteOperation(bundle, token_admin_registry.Deploy, deps, contract.DeployInput[tokenadminregistry.TokenAdminRegistry]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		ActAs:         []string{party},
		Template: tokenadminregistry.TokenAdminRegistry{
			Owner:        types.PARTY(party),
			InstanceId:   "",
			TokenConfigs: types.GENMAP{},
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "deploy TAR")
	require.NotEmpty(t, tarAddrRef.Output.Address, "TAR address")

	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("AMT"),
	}

	// Deploy pool and register with TAR in one changeset
	_, err = (DeployTokenPool{}).Apply(*env, CantonCSDeps[DeployTokenPoolConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		UserName:      user.GetId(),
		Party:         party,
		Config: DeployTokenPoolConfig{
			CcipOwner:                         party,
			PoolOwner:                         party,
			InstrumentId:                      instrumentId,
			Decimals:                          6,
			Qualifier:                         "AMT",
			TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tarAddrRef.Output.Address),
		},
	})
	require.NoError(t, err, "deploy token pool and register with TAR")

	// Verify TAR config: fetch TAR from ACS and unmarshal with UnmarshalActiveContract (uses UnmarshalCreatedEvent)
	tar, err := findTARByInstanceAddress(t.Context(), deps.StateServiceClient, party)
	require.NoError(t, err, "find TAR contract in ACS")

	// TokenConfigs is GENMAP (map); each value is a TokenConfig with optional tokenPool (PoolRegistration).
	// poolOwner is nested: config.data["tokenPool"].data["poolOwner"] (or tokenPool may be under "value" if optional).
	var found bool
	for _, v := range tar.TokenConfigs {
		configMap, ok := v.(map[string]any)
		if !ok {
			continue
		}
		// Unmarshaled form may wrap fields in "data"
		m := configMap
		if data, ok := configMap["data"].(map[string]any); ok {
			m = data
		}
		tokenPoolRaw, ok := m["tokenPool"]
		if !ok || tokenPoolRaw == nil {
			continue
		}
		// Optional: {"_type": "optional", "value": <PoolRegistration map>}
		var tokenPoolMap map[string]any
		if opt, ok := tokenPoolRaw.(map[string]any); ok {
			if val, has := opt["value"]; has && val != nil {
				tokenPoolMap, _ = val.(map[string]any)
			}
		}
		if tokenPoolMap == nil {
			tokenPoolMap, _ = tokenPoolRaw.(map[string]any)
		}
		if tokenPoolMap == nil {
			continue
		}
		// poolOwner may be under tokenPool.data when unmarshaled
		tokenPoolData := tokenPoolMap
		if data, ok := tokenPoolMap["data"].(map[string]any); ok {
			tokenPoolData = data
		}
		poolOwnerStr := optionalPartyStringFromMap(tokenPoolData, "poolOwner")
		if poolOwnerStr == "" {
			continue
		}
		require.Equal(t, party, poolOwnerStr, "TAR token config tokenPool.poolOwner should match pool owner")
		found = true

		break
	}
	require.True(t, found, "TAR should have a TokenConfig with tokenPool.poolOwner set (pool registered)")
}

// optionalPartyStringFromMap gets an optional party from a map: either direct string or {"value": "<party>"}.
func optionalPartyStringFromMap(m map[string]any, key string) string {
	raw, ok := m[key]
	if !ok || raw == nil {
		return ""
	}
	if s, ok := raw.(string); ok {
		return s
	}
	if opt, ok := raw.(map[string]any); ok {
		if v, ok := opt["value"]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}

	return ""
}

// findTARByInstanceAddress streams GetActiveContracts for TokenAdminRegistry and returns the ActiveContract whose instanceId matches the given instance address.
func findTARByInstanceAddress(ctx context.Context, stateClient apiv2.StateServiceClient, party string) (*tokenadminregistry.TokenAdminRegistry, error) {
	ledgerEndResp, err := stateClient.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, err
	}

	stream, err := stateClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{{
						IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
							TemplateFilter: &apiv2.TemplateFilter{
								TemplateId: &apiv2.Identifier{
									PackageId:  "#ccip-tokenadminregistry",
									ModuleName: "CCIP.TokenAdminRegistry",
									EntityName: "TokenAdminRegistry",
								},
								IncludeCreatedEventBlob: true,
							},
						},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, err
	}
	defer stream.CloseSend()
	for {
		ac, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		c, ok := ac.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}

		if c.ActiveContract.GetCreatedEvent().GetTemplateId().GetEntityName() != bindings.GetEntityName(tokenadminregistry.TokenAdminRegistry{}.GetTemplateID()) {
			continue
		}
		tar, err := bindings.UnmarshalActiveContract[tokenadminregistry.TokenAdminRegistry](c)
		if err != nil {
			return nil, err
		}
		if tar != nil {
			return tar, nil
		}
	}

	return nil, errors.New("no active TAR contract found for instance address")
}
