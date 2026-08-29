package ledger

import (
	"context"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

// Client submits ledger commands for a Canton participant.
type Client interface {
	Participant() canton.Participant
	ForParty(party string) canton.Participant
	SubmitCreate(ctx context.Context, actAs string, payload any) (*apiv2.SubmitAndWaitForTransactionResponse, error)
	SubmitExercise(ctx context.Context, actAs string, template interface{ GetTemplateID() string }, contractID, choice string, choiceArg any) (*apiv2.SubmitAndWaitForTransactionResponse, error)
	SubmitExerciseDisclosed(ctx context.Context, actAs string, template interface{ GetTemplateID() string }, contractID, choice string, choiceArg any, disclosed *apiv2.DisclosedContract) (*apiv2.SubmitAndWaitForTransactionResponse, error)
	SubmitExerciseMulti(ctx context.Context, actAs []string, template interface{ GetTemplateID() string }, contractID, choice string, choiceArg any, disclosed []*apiv2.DisclosedContract) (*apiv2.SubmitAndWaitForTransactionResponse, error)
}
