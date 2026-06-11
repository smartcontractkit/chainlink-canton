package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestIsKmsPublicKeyConflict(t *testing.T) {
	t.Parallel()

	conflict := status.Error(codes.Internal, "Existing public key for 12201996b8ff... is different than inserted key")
	wrapped := fmt.Errorf("RegisterKmsSigningKey: %w", conflict)

	assert.True(t, IsKmsPublicKeyConflict(wrapped))
	assert.True(t, IsKmsPublicKeyConflict(errors.New("Existing public key for abc is different than inserted key")))
	assert.False(t, IsKmsPublicKeyConflict(errors.New("permission denied")))
}
