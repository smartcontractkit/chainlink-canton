package ledger

import (
	"context"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/model"
)

type CommandService interface {
	SubmitAndWait(ctx context.Context, req *model.SubmitAndWaitRequest) (*model.SubmitAndWaitResponse, error)
	// New method added to interface
	SubmitAndWaitForTransaction(ctx context.Context, req *model.SubmitAndWaitRequest) (*model.SubmitAndWaitForTransactionResponse, error)
}

type commandService struct {
	client v2.CommandServiceClient
}

func NewCommandServiceClient(conn *grpc.ClientConn) *commandService {
	client := v2.NewCommandServiceClient(conn)
	return &commandService{
		client: client,
	}
}

func (c *commandService) SubmitAndWait(ctx context.Context, req *model.SubmitAndWaitRequest) (*model.SubmitAndWaitResponse, error) {
	protoReq := &v2.SubmitAndWaitRequest{
		Commands: commandsToProto(req.Commands),
	}

	resp, err := c.client.SubmitAndWait(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.SubmitAndWaitResponse{
		UpdateID:         resp.UpdateId,
		CompletionOffset: resp.CompletionOffset,
	}, nil
}

// SubmitAndWaitForTransaction implementation
func (c *commandService) SubmitAndWaitForTransaction(ctx context.Context, req *model.SubmitAndWaitRequest) (*model.SubmitAndWaitForTransactionResponse, error) {
	// The request structure for both Wait and WaitForTransaction is identical in terms of commands
	protoReq := &v2.SubmitAndWaitForTransactionRequest{
		Commands: commandsToProto(req.Commands),
	}

	resp, err := c.client.SubmitAndWaitForTransaction(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.SubmitAndWaitForTransactionResponse{
		UpdateID:         resp.Transaction.UpdateId,
		CompletionOffset: resp.Transaction.Offset,
		Transaction:      transactionToModel(resp.Transaction),
	}, nil
}
