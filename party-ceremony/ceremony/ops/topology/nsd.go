package topology

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// ProposeNamespaceDelegationOp publishes a namespace delegation for the
// participant's namespace key to the synchronizer.
//
// Canton equivalent:
//
//	participant1.topology.namespace_delegations.propose_delegation(
//	    namespace, nsKey, DelegationRestriction.CanSignAllMappings, store = synchronizerId)
var ProposeNamespaceDelegationOp = operations.NewOperation(
	"canton-ceremony/topology/propose-nsd",
	semver.MustParse("1.0.0"),
	"Publish namespace delegation to the synchronizer",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeNSDInput) (ProposeNSDOutput, error) {
		if in.Namespace == "" {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				errors.New("propose-nsd: namespace is required"),
			)
		}

		if in.SigningKeyB64 == "" {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				errors.New("propose-nsd: signing_key_b64 is required"),
			)
		}

		ctx := b.GetContext()

		keyBytes, err := base64.StdEncoding.DecodeString(in.SigningKeyB64)
		if err != nil {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding signing key proto: %w", err),
			)
		}
		var pk cryptov30.SigningPublicKey
		if err := proto.Unmarshal(keyBytes, &pk); err != nil {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("unmarshalling signing key proto: %w", err),
			)
		}

		// Build NamespaceDelegation mapping with CanSignAllMappings restriction.
		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_NamespaceDelegation{
				NamespaceDelegation: &protov30.NamespaceDelegation{
					Namespace: in.Namespace,
					TargetKey: &pk,
					Restriction: &protov30.NamespaceDelegation_CanSignAllMappings_{
						CanSignAllMappings: &protov30.NamespaceDelegation_CanSignAllMappings{},
					},
				},
			},
		}

		_, err = deps.Client.Authorize(ctx, 1, mapping, in.SynchronizerID, true)
		if err != nil {
			if client.IsTopologyMappingAlreadyExists(err) {
				deps.Logger.Infow("Namespace delegation already exists",
					"participant", in.ParticipantID, "namespace", in.Namespace)
			} else {
				return ProposeNSDOutput{}, fmt.Errorf("proposing namespace delegation: %w", err)
			}
		} else {
			deps.Logger.Infow("Namespace delegation proposed",
				"participant", in.ParticipantID, "namespace", in.Namespace)
		}

		return ProposeNSDOutput{
			ParticipantID:      in.ParticipantID,
			DelegationProposed: true,
		}, nil
	},
)
