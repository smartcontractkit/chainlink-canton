package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsTopologyMappingAlreadyExists(t *testing.T) {
	t.Parallel()

	st := status.New(codes.AlreadyExists, "TOPOLOGY_MAPPING_ALREADY_EXISTS(10,0): duplicate")
	wrapped := fmt.Errorf("Authorize: %w", st.Err())
	require.True(t, IsTopologyMappingAlreadyExists(wrapped))

	require.False(t, IsTopologyMappingAlreadyExists(errors.New("other error")))
}
