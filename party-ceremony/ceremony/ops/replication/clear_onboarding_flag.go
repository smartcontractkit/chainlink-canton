package replication

import (
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ClearOnboardingFlagOp clears the onboarding flag for a party on the target
// participant. This is the final step after ACS import and reconnect.
//
// Canton may need time after reconnect before the flag can be safely cleared
// (up to the synchronizer's decision timeout, typically ~1 minute). The op
// polls internally for up to 2 minutes with 10-second intervals, matching
// the Canton docs pattern:
//
//	utils.retry_until_true(timeout = 2.minutes, maxWaitPeriod = 1.minutes)
var ClearOnboardingFlagOp = operations.NewOperation(
	"canton-ceremony/replication/clear-onboarding-flag",
	semver.MustParse("1.0.0"),
	"Clear party onboarding flag on target participant after ACS import",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ClearOnboardingFlagInput) (ClearOnboardingFlagOutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || in.SynchronizerID == "" {
			return ClearOnboardingFlagOutput{}, operations.NewUnrecoverableError(
				errors.New("clear-onboarding-flag: participant_id, party_id, and synchronizer_id are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ClearOnboardingFlagOutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ClearOnboardingFlagOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		// Poll for up to 2 minutes. Canton needs time after reconnect to
		// process the ACS and reach the decision deadline before the
		// onboarding flag can be safely cleared.
		const (
			pollTimeout  = 2 * time.Minute
			pollInterval = 10 * time.Second
		)
		deadline := time.Now().Add(pollTimeout)

		deps.Logger.Infow("Polling to clear onboarding flag",
			"participant", in.ParticipantID,
			"party", in.PartyID,
			"synchronizer", in.SynchronizerID,
			"begin_offset_exclusive", in.BeginOffsetExclusive,
			"poll_timeout", pollTimeout.String(),
			"poll_interval", pollInterval.String(),
		)

		for attempt := 0; ; attempt++ {
			onboarded, callErr := deps.Client.ClearPartyOnboardingFlag(
				ctx, in.PartyID, in.SynchronizerID, in.BeginOffsetExclusive,
			)
			if callErr != nil {
				return ClearOnboardingFlagOutput{}, fmt.Errorf("clearing onboarding flag: %w", callErr)
			}

			if onboarded {
				deps.Logger.Infow("Onboarding flag cleared",
					"participant", in.ParticipantID,
					"party", in.PartyID,
					"attempts", attempt+1,
				)

				return ClearOnboardingFlagOutput{
					ParticipantID: in.ParticipantID,
					Onboarded:     true,
				}, nil
			}

			if time.Now().After(deadline) {
				return ClearOnboardingFlagOutput{}, fmt.Errorf(
					"onboarding flag not cleared for party %s after %s (%d attempts)",
					in.PartyID, pollTimeout, attempt+1,
				)
			}

			deps.Logger.Infow("Onboarding flag not yet cleared, waiting...",
				"participant", in.ParticipantID,
				"party", in.PartyID,
				"attempt", attempt+1,
				"next_retry_in", pollInterval.String(),
			)

			select {
			case <-ctx.Done():
				return ClearOnboardingFlagOutput{}, ctx.Err()
			case <-time.After(pollInterval):
			}
		}
	},
)
