package ledger

import (
	"context"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/model"
)

type CommandSubmission interface {
	Submit(ctx context.Context, req *model.SubmitRequest) (*model.SubmitResponse, error)
}

type commandSubmission struct {
	client v2.CommandSubmissionServiceClient
}

func NewCommandSubmissionClient(conn *grpc.ClientConn) *commandSubmission {
	client := v2.NewCommandSubmissionServiceClient(conn)
	return &commandSubmission{
		client: client,
	}
}

func (c *commandSubmission) Submit(ctx context.Context, req *model.SubmitRequest) (*model.SubmitResponse, error) {
	protoReq := &v2.SubmitRequest{
		Commands: commandsToProto(req.Commands),
	}

	_, err := c.client.Submit(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.SubmitResponse{}, nil
}
