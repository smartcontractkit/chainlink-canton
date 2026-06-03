package client

import (
	"context"
	"io"
	"net"
	"testing"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

func TestSetImportAcsFirstChunkMetadata(t *testing.T) {
	t.Parallel()

	const syncID = "synchronizer::1220abc"

	req := &participantv30.ImportAcsRequest{AcsSnapshot: []byte("snap")}
	require.NoError(t, setImportAcsFirstChunkMetadata(req, syncID))

	assert.Equal(t, participantv30.ContractImportMode_CONTRACT_IMPORT_MODE_ACCEPT, req.ContractImportMode)

	got, ok := importAcsSynchronizerID(req)
	require.True(t, ok, "synchronizer_id should be readable from unknown fields")
	assert.Equal(t, syncID, got)

	require.Error(t, setImportAcsFirstChunkMetadata(&participantv30.ImportAcsRequest{}, ""))
}

type captureImportAcsServer struct {
	participantv30.UnimplementedParticipantRepairServiceServer

	first *participantv30.ImportAcsRequest
}

func (s *captureImportAcsServer) ImportAcs(stream grpc.ClientStreamingServer[participantv30.ImportAcsRequest, participantv30.ImportAcsResponse]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if s.first == nil {
			s.first = proto.Clone(req).(*participantv30.ImportAcsRequest)
		}
	}

	return stream.SendAndClose(&participantv30.ImportAcsResponse{})
}

func TestGRPCClientImportAcsFirstChunk(t *testing.T) {
	t.Parallel()

	const syncID = "synchronizer::1220deadbeef"
	snapshot := []byte{0x1f, 0x8b, 0x08} // minimal gzip header prefix

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	capture := &captureImportAcsServer{}
	participantv30.RegisterParticipantRepairServiceServer(srv, capture)

	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := NewGRPCClient(conn)
	require.NoError(t, client.ImportAcs(t.Context(), snapshot, syncID))

	require.NotNil(t, capture.first, "server should receive at least one ImportAcsRequest")
	assert.Equal(t, participantv30.ContractImportMode_CONTRACT_IMPORT_MODE_ACCEPT, capture.first.ContractImportMode)
	assert.Equal(t, snapshot, capture.first.AcsSnapshot)

	got, ok := importAcsSynchronizerID(capture.first)
	require.True(t, ok, "first chunk must include synchronizer_id")
	assert.Equal(t, syncID, got)
}

func TestGRPCClientImportAcsRequiresSynchronizerID(t *testing.T) {
	t.Parallel()

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	participantv30.RegisterParticipantRepairServiceServer(srv, &captureImportAcsServer{})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	err = NewGRPCClient(conn).ImportAcs(t.Context(), []byte("x"), "")
	require.ErrorContains(t, err, "synchronizer_id is required")
}
