package usbwallet

import (
	"testing"

	interactivepb "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	devicepb "github.com/smartcontractkit/chainlink-canton/examples/cli/ledger/proto"
)

func exerciseNode(nodeID string, children ...string) *interactivepb.DamlTransaction_Node {
	return &interactivepb.DamlTransaction_Node{
		NodeId: nodeID,
		VersionedNode: &interactivepb.DamlTransaction_Node_V1{V1: &v1.Node{
			NodeType: &v1.Node_Exercise{Exercise: &v1.Exercise{Children: children}},
		}},
	}
}

func createNode(nodeID string) *interactivepb.DamlTransaction_Node {
	return &interactivepb.DamlTransaction_Node{
		NodeId: nodeID,
		VersionedNode: &interactivepb.DamlTransaction_Node_V1{V1: &v1.Node{
			NodeType: &v1.Node_Create{Create: &v1.Create{ContractId: nodeID}},
		}},
	}
}

func TestOrderNodes(t *testing.T) {
	t.Parallel()

	// Two trees, as in the example of the Ledger SPLIT_TRANSACTION documentation:
	//   0 -> {1 -> {3, 4}, 2}   and   5 -> {6, 7 -> {8, 9}}
	nodes := []*interactivepb.DamlTransaction_Node{
		exerciseNode("0", "1", "2"),
		exerciseNode("1", "3", "4"),
		createNode("2"),
		createNode("3"),
		createNode("4"),
		exerciseNode("5", "6", "7"),
		createNode("6"),
		exerciseNode("7", "8", "9"),
		createNode("8"),
		createNode("9"),
	}

	ordered, err := orderNodes([]string{"0", "5"}, nodes)
	require.NoError(t, err)

	got := make([]string, 0, len(ordered))
	for _, node := range ordered {
		got = append(got, node.GetNodeId())
	}
	require.Equal(t, []string{"3", "4", "1", "2", "0", "6", "8", "9", "7", "5"}, got)

	// Every child must be transmitted before its parent.
	position := make(map[string]int, len(got))
	for i, nodeID := range got {
		position[nodeID] = i
	}
	for _, node := range nodes {
		for _, child := range nodeChildren(node) {
			require.Less(t, position[child], position[node.GetNodeId()],
				"child %s must be sent before parent %s", child, node.GetNodeId())
		}
	}
}

func TestOrderNodesRejectsUnreachableNode(t *testing.T) {
	t.Parallel()

	nodes := []*interactivepb.DamlTransaction_Node{createNode("0"), createNode("1")}
	_, err := orderNodes([]string{"0"}, nodes)
	require.ErrorContains(t, err, "not reachable")
}

func TestOrderNodesRejectsUnknownNode(t *testing.T) {
	t.Parallel()

	_, err := orderNodes([]string{"0"}, []*interactivepb.DamlTransaction_Node{exerciseNode("0", "1")})
	require.ErrorContains(t, err, "unknown node")
}

func TestSplitPreparedTransaction(t *testing.T) {
	t.Parallel()

	transaction := &interactivepb.PreparedTransaction{
		Transaction: &interactivepb.DamlTransaction{
			Version: "2.1",
			Roots:   []string{"0"},
			Nodes:   []*interactivepb.DamlTransaction_Node{exerciseNode("0", "1"), createNode("1")},
			NodeSeeds: []*interactivepb.DamlTransaction_NodeSeed{
				{NodeId: 0, Seed: []byte{0xaa}},
			},
		},
		Metadata: &interactivepb.Metadata{
			SubmitterInfo:   &interactivepb.Metadata_SubmitterInfo{ActAs: []string{"alice"}, CommandId: "cmd"},
			SynchronizerId:  "sync",
			MediatorGroup:   1,
			TransactionUuid: "uuid",
			PreparationTime: 42,
			InputContracts: []*interactivepb.Metadata_InputContract{{
				Contract:  &interactivepb.Metadata_InputContract_V1{V1: &v1.Create{ContractId: "contract"}},
				CreatedAt: 7,
				EventBlob: []byte("this blob must not be sent to the device"),
			}},
		},
	}

	messages, err := splitPreparedTransaction(transaction)
	require.NoError(t, err)
	// DAML transaction + 2 nodes + metadata + 1 input contract.
	require.Len(t, messages, 5)

	var damlTransaction devicepb.DeviceDamlTransaction
	require.NoError(t, proto.Unmarshal(messages[0], &damlTransaction))
	require.Equal(t, "2.1", damlTransaction.GetVersion())
	require.Equal(t, []string{"0"}, damlTransaction.GetRoots())
	require.Equal(t, int32(2), damlTransaction.GetNodesCount())
	require.Len(t, damlTransaction.GetNodeSeeds(), 1)

	// The child create node is sent before its exercise parent.
	var firstNode, secondNode devicepb.DeviceDamlTransaction_Node
	require.NoError(t, proto.Unmarshal(messages[1], &firstNode))
	require.NoError(t, proto.Unmarshal(messages[2], &secondNode))
	require.Equal(t, "1", firstNode.GetNodeId())
	require.Equal(t, "0", secondNode.GetNodeId())

	var metadata devicepb.DeviceMetadata
	require.NoError(t, proto.Unmarshal(messages[3], &metadata))
	require.Equal(t, "sync", metadata.GetSynchronizerId())
	require.Equal(t, int32(1), metadata.GetInputContractsCount())
	require.Equal(t, []string{"alice"}, metadata.GetSubmitterInfo().GetActAs())

	var inputContract devicepb.DeviceMetadata_InputContract
	require.NoError(t, proto.Unmarshal(messages[4], &inputContract))
	require.Equal(t, uint64(7), inputContract.GetCreatedAt())
	require.Equal(t, "contract", inputContract.GetV1().GetContractId())
	require.NotContains(t, string(messages[4]), "this blob must not be sent to the device")
}

func TestSplitPreparedTransactionPreservesDriverMetadata(t *testing.T) {
	t.Parallel()

	driverMetadata := []byte("driver metadata")

	inputContract := &interactivepb.Metadata_InputContract{
		Contract:  &interactivepb.Metadata_InputContract_V1{V1: &v1.Create{ContractId: "contract"}},
		CreatedAt: 7,
	}
	// The generated dazl client has no driver_metadata field, so simulate a participant that
	// sends one by appending it as an unknown field.
	unknown := protowire.AppendTag(nil, driverMetadataField, protowire.BytesType)
	unknown = protowire.AppendBytes(unknown, driverMetadata)
	inputContract.ProtoReflect().SetUnknown(unknown)

	encoded, err := convertInputContract(inputContract)
	require.NoError(t, err)

	var deviceContract devicepb.DeviceMetadata_InputContract
	require.NoError(t, proto.Unmarshal(encoded, &deviceContract))
	require.Equal(t, driverMetadata, deviceContract.GetDriverMetadata())
}
