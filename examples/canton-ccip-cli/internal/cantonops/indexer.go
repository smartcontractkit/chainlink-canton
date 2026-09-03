package cantonops

import (
	"context"
	"fmt"
	"net/http"
	"time"

	v1 "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/api/handlers/v1"
	indexerclient "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/client"
)

// WaitForVerifierResult polls the indexer for verifier results for the given
// message id until at least one is available, the context is cancelled or the
// timeout elapses.
func WaitForVerifierResult(
	ctx context.Context,
	idx *indexerclient.IndexerClient,
	messageID string,
	timeout time.Duration,
) (v1.VerifierResultsByMessageIDResponse, error) {
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return v1.VerifierResultsByMessageIDResponse{}, fmt.Errorf("timed out waiting for verifier results for message %s", messageID)
		}
		status, resp, err := idx.VerifierResultsByMessageID(ctx, v1.VerifierResultsByMessageIDInput{MessageID: messageID})
		if err != nil {
			fmt.Printf("indexer query error (retrying): %v\n", err)
		} else if status != http.StatusOK {
			fmt.Printf("indexer status %d (retrying)\n", status)
		} else if len(resp.Results) > 0 {
			return resp, nil
		}
		select {
		case <-ctx.Done():
			return v1.VerifierResultsByMessageIDResponse{}, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}
