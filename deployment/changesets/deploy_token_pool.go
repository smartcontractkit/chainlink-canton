package changesets

import (
	"fmt"

	"github.com/aws/smithy-go/ptr"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/auth"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// DeployTokenPoolConfig is the config for deploying a LockReleaseTokenPool.
type DeployTokenPoolConfig struct {
	CcipOwner    string
	PoolOwner    string
	InstrumentId lockreleasetokenpool.InstrumentId
	Decimals     int64
	// Qualifier is optional (e.g. token symbol) for AddressRef and idempotency.
	Qualifier string
	// Optional; defaults to empty. ChainCCVRequirements can be set for chain-specific CCV requirements.
	ChainCCVRequirements types.GENMAP
	// Optional; defaults to empty. PoolReceiveContext can be set for receive context.
	PoolReceiveContext lockreleasetokenpool.ChoiceContext
	// Optional; defaults to 24h RelativeHours. TransferTimeout for the pool.
	TransferTimeout lockreleasetokenpool.TransferTimeout
}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployTokenPoolConfig]] = DeployTokenPool{}

type DeployTokenPool struct{}

func (d DeployTokenPool) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployTokenPoolConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}
	return nil
}

func (d DeployTokenPool) Apply(e cldf.Environment, config CantonCSDeps[DeployTokenPoolConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	token, err := participant.JWTProvider.Token(e.GetContext())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get JWT: %w", err)
	}

	insecureCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	ledgerApiClient, err := grpc.NewClient(participant.Endpoints.GRPCLedgerAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create gRPC ledger API client: %w", err)
	}

	deps := dependencies.CantonDeps{
		Chain:                chain,
		CommandServiceClient: apiv2.NewCommandServiceClient(ledgerApiClient),
		StateServiceClient:   apiv2.NewStateServiceClient(ledgerApiClient),
		Party:                config.Party,
	}

	cfg := config.Config
	chainCCVReqs := cfg.ChainCCVRequirements
	if chainCCVReqs == nil {
		chainCCVReqs = types.GENMAP{}
	}
	poolReceiveContext := cfg.PoolReceiveContext
	if poolReceiveContext.Values == nil {
		poolReceiveContext = lockreleasetokenpool.ChoiceContext{Values: types.TEXTMAP{}}
	}
	transferTimeout := cfg.TransferTimeout
	if transferTimeout.RelativeHours == nil && transferTimeout.Indefinite == nil {
		h := types.INT64(24)
		transferTimeout = lockreleasetokenpool.TransferTimeout{RelativeHours: &h}
	}

	template := lockreleasetokenpool.LockReleaseTokenPool{
		CcipOwner:            types.PARTY(cfg.CcipOwner),
		PoolOwner:            types.PARTY(cfg.PoolOwner),
		InstanceId:           "", // set by deploy operation
		InstrumentId:         cfg.InstrumentId,
		Decimals:             types.INT64(cfg.Decimals),
		ChainCCVRequirements: chainCCVReqs,
		PoolReceiveContext:   poolReceiveContext,
		TransferTimeout:      transferTimeout,
	}

	qualifier := ptr.String(cfg.Qualifier)
	if cfg.Qualifier == "" {
		qualifier = nil
	}

	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, lock_release_token_pool.Deploy, deps, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
		ChainSelector: config.ChainSelector,
		Qualifier:     qualifier,
		ActAs:         []string{cfg.PoolOwner},
		Template:      template,
		OwnerParty:    types.PARTY(cfg.PoolOwner),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply DeployTokenPool operation: %w", err)
	}

	if err = ds.AddressRefStore.Add(out.Output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save deployed LockReleaseTokenPool contract address: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []cld_ops.Report[any, any]{},
	}, nil
}
