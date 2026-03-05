package changesets

import (
	"fmt"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type UnvetAndReuploadDARsConfig struct {
	// DAR payloads to upload after unvetting.
	DARs [][]byte
	// Optional explicit list of main package IDs to unvet first.
	// If empty, all currently listed DARs on the participant are unvetted.
	MainPackageIDsToUnvet []string
	// Whether to wait for package vetting changes to synchronize.
	SynchronizeVetting bool
	// Optional synchronizer ID for vet/unvet requests.
	SynchronizerID *string
}

var _ cldf.ChangeSetV2[CantonCSDeps[UnvetAndReuploadDARsConfig]] = UnvetAndReuploadDARs{}

type UnvetAndReuploadDARs struct{}

func (u UnvetAndReuploadDARs) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[UnvetAndReuploadDARsConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf(
			"participant index %d out of range for canton chain %d with %d participants",
			config.Participant,
			config.ChainSelector,
			len(chain.Participants),
		)
	}
	if len(config.Config.DARs) == 0 {
		return fmt.Errorf("at least one DAR is required for reupload")
	}

	return nil
}

func (u UnvetAndReuploadDARs) Apply(e cldf.Environment, config CantonCSDeps[UnvetAndReuploadDARsConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	mainPackageIDs, err := mainPackageIDsToUnvet(e, participant, config.Config)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for _, mainPackageID := range mainPackageIDs {
		_, err = participant.AdminServices.Package.UnvetDar(e.GetContext(), &participantv30.UnvetDarRequest{
			MainPackageId:  mainPackageID,
			SynchronizerId: config.Config.SynchronizerID,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to unvet DAR %q: %w", mainPackageID, err)
		}
	}

	darData := make([]*participantv30.UploadDarRequest_UploadDarData, 0, len(config.Config.DARs))
	for i, dar := range config.Config.DARs {
		if len(dar) == 0 {
			return cldf.ChangesetOutput{}, fmt.Errorf("DAR at index %d is empty", i)
		}
		darData = append(darData, &participantv30.UploadDarRequest_UploadDarData{
			Bytes: dar,
		})
	}

	_, err = participant.AdminServices.Package.UploadDar(e.GetContext(), &participantv30.UploadDarRequest{
		Dars:               darData,
		VetAllPackages:     true,
		SynchronizeVetting: config.Config.SynchronizeVetting,
		SynchronizerId:     config.Config.SynchronizerID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to reupload DARs: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []operations.Report[any, any]{},
	}, nil
}

func mainPackageIDsToUnvet(e cldf.Environment, participant canton.Participant, cfg UnvetAndReuploadDARsConfig) ([]string, error) {
	if len(cfg.MainPackageIDsToUnvet) > 0 {
		return cfg.MainPackageIDsToUnvet, nil
	}

	listResp, err := participant.AdminServices.Package.ListDars(e.GetContext(), &participantv30.ListDarsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list existing DARs before unvet: %w", err)
	}

	mainPackageIDs := make([]string, 0, len(listResp.GetDars()))
	seen := make(map[string]struct{}, len(listResp.GetDars()))
	for _, dar := range listResp.GetDars() {
		main := dar.GetMain()
		if main == "" {
			continue
		}
		if _, ok := seen[main]; ok {
			continue
		}
		seen[main] = struct{}{}
		mainPackageIDs = append(mainPackageIDs, main)
	}

	return mainPackageIDs, nil
}
