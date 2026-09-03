package usbwallet

import (
	"errors"
	"fmt"

	interactivepb "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	devicepb "github.com/smartcontractkit/chainlink-canton/examples/canton-ccip-cli/ledger/proto"
)

// driverMetadataField is the field number of Metadata.InputContract.driver_metadata. The
// generated dazl client does not know the field yet, so it is recovered from the unknown fields
// to keep the device side hash identical to the one computed by the participant.
const driverMetadataField protowire.Number = 1001

// deviceMessage is one of the component messages a prepared transaction is split into. The label
// is only used to point at the offending component when the device rejects one of them.
type deviceMessage struct {
	label   string
	payload []byte
}

// splitPreparedTransaction serializes a prepared transaction into the component messages the
// device expects, in transmission order: the DAML transaction, every node (children before
// their parents), the metadata and finally every input contract.
func splitPreparedTransaction(transaction *interactivepb.PreparedTransaction) ([]deviceMessage, error) {
	damlTransaction := transaction.GetTransaction()
	if damlTransaction == nil {
		return nil, errors.New("ledger: prepared transaction has no DAML transaction")
	}
	metadata := transaction.GetMetadata()
	if metadata == nil {
		return nil, errors.New("ledger: prepared transaction has no metadata")
	}

	nodeSeeds := make([]*devicepb.DeviceDamlTransaction_NodeSeed, 0, len(damlTransaction.GetNodeSeeds()))
	for _, seed := range damlTransaction.GetNodeSeeds() {
		nodeSeeds = append(nodeSeeds, &devicepb.DeviceDamlTransaction_NodeSeed{
			NodeId: seed.GetNodeId(),
			Seed:   seed.GetSeed(),
		})
	}

	// The nodes are stripped from the DAML transaction and replaced by their count, since the
	// device does not have enough memory to hold the whole transaction at once.
	deviceTransaction, err := proto.Marshal(&devicepb.DeviceDamlTransaction{
		Version:    damlTransaction.GetVersion(),
		Roots:      damlTransaction.GetRoots(),
		NodesCount: int32(len(damlTransaction.GetNodes())), //nolint:gosec
		NodeSeeds:  nodeSeeds,
	})
	if err != nil {
		return nil, fmt.Errorf("ledger: failed to marshal DAML transaction: %w", err)
	}

	nodes, err := orderNodes(damlTransaction.GetRoots(), damlTransaction.GetNodes())
	if err != nil {
		return nil, err
	}

	messages := make([]deviceMessage, 0, 2+len(nodes)+len(metadata.GetInputContracts()))
	messages = append(messages, deviceMessage{label: "DAML transaction", payload: deviceTransaction})

	for i, node := range nodes {
		deviceNode := &devicepb.DeviceDamlTransaction_Node{NodeId: node.GetNodeId()}
		if versioned := node.GetV1(); versioned != nil {
			deviceNode.VersionedNode = &devicepb.DeviceDamlTransaction_Node_V1{V1: versioned}
		}

		encoded, err := proto.Marshal(deviceNode)
		if err != nil {
			return nil, fmt.Errorf("ledger: failed to marshal node %q: %w", node.GetNodeId(), err)
		}
		messages = append(messages, deviceMessage{
			label:   fmt.Sprintf("node %d/%d (node id %q)", i+1, len(nodes), node.GetNodeId()),
			payload: encoded,
		})
	}

	deviceMetadata, err := convertMetadata(metadata)
	if err != nil {
		return nil, err
	}
	messages = append(messages, deviceMessage{label: "metadata", payload: deviceMetadata})

	for i, inputContract := range metadata.GetInputContracts() {
		encoded, err := convertInputContract(inputContract)
		if err != nil {
			return nil, fmt.Errorf("ledger: failed to marshal input contract %d: %w", i, err)
		}
		messages = append(messages, deviceMessage{
			label:   fmt.Sprintf("input contract %d/%d", i+1, len(metadata.GetInputContracts())),
			payload: encoded,
		})
	}

	return messages, nil
}

func convertMetadata(metadata *interactivepb.Metadata) ([]byte, error) {
	deviceMetadata := &devicepb.DeviceMetadata{
		SynchronizerId:      metadata.GetSynchronizerId(),
		MediatorGroup:       metadata.GetMediatorGroup(),
		TransactionUuid:     metadata.GetTransactionUuid(),
		PreparationTime:     metadata.GetPreparationTime(),
		InputContractsCount: int32(len(metadata.GetInputContracts())), //nolint:gosec
	}
	if submitterInfo := metadata.GetSubmitterInfo(); submitterInfo != nil {
		deviceMetadata.SubmitterInfo = &devicepb.DeviceMetadata_SubmitterInfo{
			ActAs:     submitterInfo.GetActAs(),
			CommandId: submitterInfo.GetCommandId(),
		}
	}
	if metadata.MinLedgerEffectiveTime != nil {
		deviceMetadata.MinLedgerEffectiveTime = new(metadata.GetMinLedgerEffectiveTime())
	}
	if metadata.MaxLedgerEffectiveTime != nil {
		deviceMetadata.MaxLedgerEffectiveTime = new(metadata.GetMaxLedgerEffectiveTime())
	}
	// DeviceMetadata deliberately has no global_key_mapping field, the device does not hash it.

	encoded, err := proto.Marshal(deviceMetadata)
	if err != nil {
		return nil, fmt.Errorf("ledger: failed to marshal metadata: %w", err)
	}

	return encoded, nil
}

// convertInputContract drops the event blob, which is not part of the hash and would only waste
// APDU bandwidth.
func convertInputContract(inputContract *interactivepb.Metadata_InputContract) ([]byte, error) {
	deviceContract := &devicepb.DeviceMetadata_InputContract{
		CreatedAt:      inputContract.GetCreatedAt(),
		DriverMetadata: unknownBytesField(inputContract, driverMetadataField),
	}
	if create := inputContract.GetV1(); create != nil {
		deviceContract.Contract = &devicepb.DeviceMetadata_InputContract_V1{V1: create}
	}

	return proto.Marshal(deviceContract)
}

// unknownBytesField returns the value of a length delimited field that the generated message does
// not know about, or nil if the field is absent.
func unknownBytesField(message proto.Message, field protowire.Number) []byte {
	unknown := message.ProtoReflect().GetUnknown()
	for len(unknown) > 0 {
		number, typ, n := protowire.ConsumeTag(unknown)
		if n < 0 {
			return nil
		}
		unknown = unknown[n:]

		if number == field && typ == protowire.BytesType {
			value, n := protowire.ConsumeBytes(unknown)
			if n < 0 {
				return nil
			}

			return value
		}

		n = protowire.ConsumeFieldValue(number, typ, unknown)
		if n < 0 {
			return nil
		}
		unknown = unknown[n:]
	}

	return nil
}

// orderNodes returns the transaction nodes in the order the device expects them: the nodes of
// each root tree in post order, so that the children of a node are always transmitted before the
// node itself, and the root trees in the order given by roots.
func orderNodes(roots []string, nodes []*interactivepb.DamlTransaction_Node) ([]*interactivepb.DamlTransaction_Node, error) {
	byID := make(map[string]*interactivepb.DamlTransaction_Node, len(nodes))
	for _, node := range nodes {
		if _, duplicate := byID[node.GetNodeId()]; duplicate {
			return nil, fmt.Errorf("ledger: duplicate node id %q", node.GetNodeId())
		}
		byID[node.GetNodeId()] = node
	}

	ordered := make([]*interactivepb.DamlTransaction_Node, 0, len(nodes))
	visited := make(map[string]bool, len(nodes))

	var visit func(nodeID string) error
	visit = func(nodeID string) error {
		if visited[nodeID] {
			return fmt.Errorf("ledger: node %q is reachable more than once", nodeID)
		}
		visited[nodeID] = true

		node, ok := byID[nodeID]
		if !ok {
			return fmt.Errorf("ledger: transaction references unknown node %q", nodeID)
		}
		for _, child := range nodeChildren(node) {
			if err := visit(child); err != nil {
				return err
			}
		}
		ordered = append(ordered, node)

		return nil
	}

	for _, root := range roots {
		if err := visit(root); err != nil {
			return nil, err
		}
	}

	if len(ordered) != len(nodes) {
		return nil, fmt.Errorf("ledger: %d of %d nodes are not reachable from the transaction roots", len(nodes)-len(ordered), len(nodes))
	}

	return ordered, nil
}

// nodeChildren returns the node IDs of the children of a node. Only exercise and rollback nodes
// have children.
func nodeChildren(node *interactivepb.DamlTransaction_Node) []string {
	versioned := node.GetV1()
	if versioned == nil {
		return nil
	}

	switch nodeType := versioned.GetNodeType().(type) {
	case *v1.Node_Exercise:
		return nodeType.Exercise.GetChildren()
	case *v1.Node_Rollback:
		return nodeType.Rollback.GetChildren()
	default:
		return nil
	}
}
