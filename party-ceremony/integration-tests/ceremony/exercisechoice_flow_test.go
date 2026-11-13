package tests

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/smartcontractkit/go-daml/pkg/model"
	daml "github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	registryapp "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/utility/registry_app_v0"
	registryv0 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/utility/registry_v0"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/exercisechoice"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

const linkInstrumentID = "LINK"

type ExerciseChoiceFlowTestSuite struct {
	OnboardingFlowTestSuite
}

func (s *ExerciseChoiceFlowTestSuite) SetupSuite() {
	s.OnboardingFlowTestSuite.SetupSuite()
	s.performOnboarding(s.T(), operations.NewMemoryReporter())
}

// TestOnboardingFlow overrides the inherited test; onboarding ran in SetupSuite.
func (s *ExerciseChoiceFlowTestSuite) TestOnboardingFlow() {
	s.Require().NotEmpty(s.PartyID, "PartyID should have been set during SetupSuite")
}

// TestRegistrySetupFlow sets up the LINK Registry with the decentralized party
// as provider and registrar. LocalNet has no DA operator, so the party is also
// the operator and accepts its own provider request.
func (s *ExerciseChoiceFlowTestSuite) TestRegistrySetupFlow() {
	t := s.T()
	party := types.PARTY(s.PartyID)

	kmsCfgs := []client.KMSConfig{
		s.kmsConfigFor(0, "onboarding"),
		s.kmsConfigFor(1, "onboarding"),
		s.kmsConfigFor(2, "onboarding"),
	}
	deps, recorders := s.ledgerCeremonyDeps(t, "registry-setup", kmsCfgs, utilityDARLoader)

	opCfg := s.createAsParty(t, deps, registryapp.OperatorConfiguration{Operator: party})
	providerReq := s.createAsParty(t, deps, registryapp.ProviderServiceRequest{Operator: party, Provider: party})

	providerSvc := s.exerciseAsParty(t, deps, registryapp.ProviderServiceRequest{}, providerReq,
		"ProviderServiceRequest_Accept",
		registryapp.ProviderServiceRequestAccept{OperatorConfigurationCid: types.CONTRACT_ID(opCfg)},
	).created(t, "ProviderService")

	providerCfg := s.exerciseAsParty(t, deps, registryapp.ProviderService{}, providerSvc,
		"ProviderService_CreateProviderConfiguration",
		registryapp.ProviderServiceCreateProviderConfiguration{},
	).created(t, "ProviderConfiguration")

	createTransferRule, createAllocationFactory := types.BOOL(true), types.BOOL(true)
	registrarReq := s.createAsParty(t, deps, registryapp.RegistrarServiceRequest{
		Operator:                party,
		Provider:                party,
		Registrar:               party,
		CreateTransferRule:      &createTransferRule,
		CreateAllocationFactory: &createAllocationFactory,
	})

	accepted := s.exerciseAsParty(t, deps, registryapp.ProviderService{}, providerSvc,
		"ProviderService_AcceptRegistrarServiceRequest",
		registryapp.ProviderServiceAcceptRegistrarServiceRequest{
			Cid:     types.CONTRACT_ID(registrarReq),
			Payload: registryapp.RegistrarServiceRequestAccept{ProviderConfigurationCid: types.CONTRACT_ID(providerCfg)},
		},
	)
	registrarSvc := accepted.created(t, "RegistrarService")
	accepted.created(t, "AllocationFactory")

	instrumentCfg := s.exerciseAsParty(t, deps, registryapp.RegistrarService{}, registrarSvc,
		"RegistrarService_CreateInstrumentConfiguration",
		registryapp.RegistrarServiceCreateInstrumentConfiguration{InstrumentId: linkInstrumentID},
	).created(t, "InstrumentConfiguration")

	_, module, entity := splitTemplateID(t, registryv0.InstrumentConfiguration{}.GetTemplateID())
	active, err := deps[0].LedgerClient.GetActiveContractsByTemplateForParty(t.Context(), s.PartyID, registryv0.PackageName, module, entity)
	require.NoError(t, err, "query InstrumentConfiguration")
	require.Len(t, active, 1, "one InstrumentConfiguration is active for the party")
	require.Equal(t, instrumentCfg, active[0].GetContractId())
	require.Contains(t, protojson.Format(active[0].GetCreateArguments()), `"`+linkInstrumentID+`"`, "the configuration is for LINK")

	out := s.runQueryContractsCLI(t, registryv0.InstrumentConfiguration{}.GetTemplateID())
	require.Contains(t, out, "Active contracts: 1\n", "query-contracts finds exactly the configuration")
	require.Contains(t, out, "Contract ID: "+instrumentCfg)
	require.Contains(t, out, `"text":"`+linkInstrumentID+`"`)

	assertRecordersMatchKMSConfig(t, kmsCfgs, recorders)
}

func utilityDARLoader(name, version string) ([]byte, error) {
	return contracts.GetDar(contracts.Package(name), version)
}

// utilityPackages puts the Registry app first because the contract deploy
// ceremony creates templates from the first package.
var utilityPackages = []ledger.PackageRef{
	{Name: string(contracts.UtilityRegistryAppV0), Version: contracts.CurrentVersion},
	{Name: string(contracts.UtilityCommercialsV0), Version: contracts.CurrentVersion},
	{Name: string(contracts.UtilityCredentialV0), Version: contracts.CurrentVersion},
	{Name: string(contracts.UtilityCredentialAppV0), Version: contracts.CurrentVersion},
	{Name: string(contracts.UtilityRegistryV0), Version: contracts.CurrentVersion},
	{Name: string(contracts.UtilityRegistryHoldingV0), Version: contracts.CurrentVersion},
}

type damlTemplate interface {
	GetTemplateID() string
	CreateCommand() *model.CreateCommand
}

// createAsParty creates a contract through the contract deploy ceremony.
func (s *ExerciseChoiceFlowTestSuite) createAsParty(t *testing.T, deps []ledger.ContractDeployDeps, tmpl damlTemplate) string {
	t.Helper()
	_, module, entity := splitTemplateID(t, tmpl.GetTemplateID())
	args, err := protojson.Marshal(daml.MapToRecord(tmpl.CreateCommand().Arguments))
	require.NoError(t, err, "encode %s create arguments", entity)

	out := runUntilComplete(t, contractdeploy.ContractDeploySequence, deps, contractdeploy.ContractDeployInput{
		DecentralizedPartyID: s.PartyID,
		SynchronizerID:       s.SynchronizerID,
		Packages:             utilityPackages,
		TemplateModule:       module,
		TemplateEntity:       entity,
		ContractArgs:         string(args),
	}, contractdeploy.ErrThresholdNotMet)
	require.NotEmpty(t, out.ContractID, "%s created", entity)
	t.Logf("Created %s: %s", entity, out.ContractID)

	return out.ContractID
}

type exercised exercisechoice.ExerciseChoiceOutput

func (e exercised) created(t *testing.T, entity string) string {
	t.Helper()
	for _, c := range e.CreatedContracts {
		if strings.HasSuffix(c.TemplateID, ":"+entity) {
			return c.ContractID
		}
	}
	require.Failf(t, "contract not created", "update %s created no %s: %+v", e.UpdateID, entity, e.CreatedContracts)

	return ""
}

// exerciseAsParty exercises a choice through the exercise choice ceremony.
func (s *ExerciseChoiceFlowTestSuite) exerciseAsParty(
	t *testing.T,
	deps []ledger.ContractDeployDeps,
	tmpl interface{ GetTemplateID() string },
	contractID, choice string,
	arg any,
) exercised {
	t.Helper()
	pkg, module, entity := splitTemplateID(t, tmpl.GetTemplateID())
	encoded, err := protojson.Marshal(daml.MapToValue(arg))
	require.NoError(t, err, "encode %s argument", choice)

	out := runUntilComplete(t, exercisechoice.ExerciseChoiceSequence, deps, exercisechoice.ExerciseChoiceInput{
		DecentralizedPartyID: s.PartyID,
		SynchronizerID:       s.SynchronizerID,
		PackageID:            pkg,
		TemplateModule:       module,
		TemplateEntity:       entity,
		ContractID:           contractID,
		Choice:               choice,
		ChoiceArgument:       string(encoded),
	}, exercisechoice.ErrThresholdNotMet)
	require.NotEmpty(t, out.UpdateID, "%s committed", choice)
	t.Logf("Exercised %s: update %s", choice, out.UpdateID)

	return exercised(out)
}

// runUntilComplete runs the sequence as each participant in turn, as the
// participants would with `resume`.
func runUntilComplete[I, O any](
	t *testing.T,
	seq *operations.Sequence[I, O, ledger.ContractDeployDeps],
	deps []ledger.ContractDeployDeps,
	input I,
	thresholdErr error,
) O {
	t.Helper()
	reporter := operations.NewMemoryReporter()
	for run := range 3 * len(deps) {
		bundle := operations.NewBundle(t.Context, logger.Test(t), reporter)
		sr, err := operations.ExecuteSequence(bundle, seq, deps[run%len(deps)], input)
		if err == nil {
			return sr.Output
		}
		require.ErrorContains(t, err, thresholdErr.Error(), "run %d as participant %d", run+1, run%len(deps)+1)
	}
	require.FailNow(t, "sequence did not complete", "%s after %d runs", seq.ID(), 3*len(deps))

	var zero O

	return zero
}

func splitTemplateID(t *testing.T, templateID string) (pkg, module, entity string) {
	t.Helper()
	parts := strings.Split(templateID, ":")
	require.Len(t, parts, 3, "template ID %q", templateID)

	return parts[0], parts[1], parts[2]
}

// runQueryContractsCLI builds the CLI and runs query-contracts as participant 1,
// the recovery command the runbook gives operators.
func (s *ExerciseChoiceFlowTestSuite) runQueryContractsCLI(t *testing.T, templateID string) string {
	t.Helper()
	p := s.chain.Participants[0]

	bin := filepath.Join(t.TempDir(), "canton-party-ceremony")
	build := exec.CommandContext(t.Context(), "go", "build", "-o", bin, ".")
	build.Dir = filepath.Join("..", "..")
	out, err := build.CombinedOutput()
	require.NoError(t, err, "build CLI: %s", out)

	host, portText, err := net.SplitHostPort(p.Endpoints.GRPCLedgerAPIURL)
	require.NoError(t, err, "ledger API URL %q", p.Endpoints.GRPCLedgerAPIURL)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)
	cfg := map[string]any{"participant_id": p.Name, "ledger_host": host, "ledger_port": port}
	if p.TokenSource != nil {
		tok, tokErr := p.TokenSource.Token()
		require.NoError(t, tokErr)
		cfg["ledger_jwt"] = tok.AccessToken
	}
	cfgJSON, err := json.Marshal(cfg)
	require.NoError(t, err)
	cfgPath := filepath.Join(t.TempDir(), "participant-config.json")
	require.NoError(t, os.WriteFile(cfgPath, cfgJSON, 0o600))

	query := exec.CommandContext(t.Context(), bin, "query-contracts", "--party", s.PartyID, "--template-id", templateID, "--config", cfgPath)
	out, err = query.CombinedOutput()
	require.NoError(t, err, "query-contracts: %s", out)
	t.Logf("query-contracts output:\n%s", out)

	return string(out)
}
