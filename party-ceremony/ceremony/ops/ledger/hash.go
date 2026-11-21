package ledger

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
)

// HashPreparedTransactionV2 computes the hash a party signs for tx, following
// Canton's external signing hashing scheme V2:
// https://github.com/digital-asset/canton/blob/main/community/app/src/pack/examples/08-interactive-submission/daml_transaction_hashing_v2.py
func HashPreparedTransactionV2(tx *interactive.PreparedTransaction) ([]byte, error) {
	h := txHasher{
		nodes: make(map[string]*interactive.DamlTransaction_Node, len(tx.GetTransaction().GetNodes())),
		seeds: make(map[string][]byte, len(tx.GetTransaction().GetNodeSeeds())),
	}
	for _, n := range tx.GetTransaction().GetNodes() {
		h.nodes[n.GetNodeId()] = n
	}
	for _, s := range tx.GetTransaction().GetNodeSeeds() {
		h.seeds[strconv.Itoa(int(s.GetNodeId()))] = s.GetSeed()
	}

	txEnc := &encoder{}
	txEnc.string(tx.GetTransaction().GetVersion())
	txEnc.int32(int32(len(tx.GetTransaction().GetRoots())))
	for _, root := range tx.GetTransaction().GetRoots() {
		txEnc.raw(h.nodeHash(root, txEnc))
	}
	if txEnc.err != nil {
		return nil, fmt.Errorf("encoding transaction: %w", txEnc.err)
	}

	metaEnc := &encoder{}
	metaEnc.metadata(tx.GetMetadata())
	if metaEnc.err != nil {
		return nil, fmt.Errorf("encoding metadata: %w", metaEnc.err)
	}

	final := sha256.New()
	final.Write(preparedTxHashPurpose)
	final.Write([]byte{byte(interactive.HashingSchemeVersion_HASHING_SCHEME_VERSION_V2)})
	final.Write(purposeHash(txEnc.buf))
	final.Write(purposeHash(metaEnc.buf))

	return final.Sum(nil), nil
}

var preparedTxHashPurpose = []byte{0x00, 0x00, 0x00, 0x30}

const (
	nodeEncodingVersion     = 0x01
	metadataEncodingVersion = 0x01
)

func purposeHash(encoded []byte) []byte {
	sum := sha256.Sum256(append(append([]byte{}, preparedTxHashPurpose...), encoded...))

	return sum[:]
}

type txHasher struct {
	nodes map[string]*interactive.DamlTransaction_Node
	seeds map[string][]byte
}

// nodeHash returns the hash of the node with the given ID, recording any
// failure on parent so the whole encoding fails.
func (h txHasher) nodeHash(id string, parent *encoder) []byte {
	node, ok := h.nodes[id]
	if !ok {
		parent.fail(fmt.Errorf("node %q is not in the transaction", id))

		return nil
	}
	v1Node := node.GetV1()
	if v1Node == nil {
		parent.fail(fmt.Errorf("node %q has an unsupported version", id))

		return nil
	}

	e := &encoder{}
	switch n := v1Node.GetNodeType().(type) {
	case *v1.Node_Create:
		seed, hasSeed := h.seeds[id]
		e.create(n.Create, seed, hasSeed)
	case *v1.Node_Exercise:
		seed, hasSeed := h.seeds[id]
		if !hasSeed {
			parent.fail(fmt.Errorf("exercise node %q has no seed", id))

			return nil
		}
		h.exercise(e, n.Exercise, seed)
	case *v1.Node_Fetch:
		e.fetch(n.Fetch)
	case *v1.Node_Rollback:
		e.byte(nodeEncodingVersion)
		e.byte(0x03)
		h.children(e, n.Rollback.GetChildren())
	default:
		parent.fail(fmt.Errorf("node %q has an unsupported type", id))

		return nil
	}
	if e.err != nil {
		parent.fail(fmt.Errorf("node %q: %w", id, e.err))

		return nil
	}
	sum := sha256.Sum256(e.buf)

	return sum[:]
}

func (h txHasher) children(e *encoder, ids []string) {
	e.int32(int32(len(ids)))
	for _, id := range ids {
		e.raw(h.nodeHash(id, e))
	}
}

func (h txHasher) exercise(e *encoder, x *v1.Exercise, seed []byte) {
	e.byte(nodeEncodingVersion)
	e.string(x.GetLfVersion())
	e.byte(0x01)
	e.raw(seed)
	e.hexString(x.GetContractId())
	e.string(x.GetPackageName())
	e.identifier(x.GetTemplateId())
	e.strings(x.GetSignatories())
	e.strings(x.GetStakeholders())
	e.strings(x.GetActingParties())
	if e.present(x.GetInterfaceId() != nil) {
		e.identifier(x.GetInterfaceId())
	}
	e.string(x.GetChoiceId())
	e.value(x.GetChosenValue())
	e.bool(x.GetConsuming())
	if e.present(x.GetExerciseResult() != nil) {
		e.value(x.GetExerciseResult())
	}
	e.strings(x.GetChoiceObservers())
	h.children(e, x.GetChildren())
}

// encoder appends the V2 encoding of each value to buf. The first failure is
// kept in err and later writes still run, so callers check err once.
type encoder struct {
	buf []byte
	err error
}

func (e *encoder) fail(err error) {
	if e.err == nil {
		e.err = err
	}
}

func (e *encoder) raw(b []byte) { e.buf = append(e.buf, b...) }

func (e *encoder) byte(b byte) { e.buf = append(e.buf, b) }

func (e *encoder) bool(v bool) {
	if v {
		e.byte(0x01)
	} else {
		e.byte(0x00)
	}
}

// present writes the optional marker and reports whether the value follows.
func (e *encoder) present(set bool) bool {
	e.bool(set)

	return set
}

func (e *encoder) int32(v int32) { e.buf = binary.BigEndian.AppendUint32(e.buf, uint32(v)) }

func (e *encoder) int64(v int64) { e.buf = binary.BigEndian.AppendUint64(e.buf, uint64(v)) }

func (e *encoder) bytes(b []byte) {
	e.int32(int32(len(b)))
	e.raw(b)
}

func (e *encoder) string(s string) { e.bytes([]byte(s)) }

func (e *encoder) hexString(s string) {
	b, err := hex.DecodeString(s)
	if err != nil {
		e.fail(fmt.Errorf("contract ID %q is not hex: %w", s, err))
	}
	e.bytes(b)
}

func (e *encoder) strings(ss []string) {
	e.int32(int32(len(ss)))
	for _, s := range ss {
		e.string(s)
	}
}

func (e *encoder) identifier(id *apiv2.Identifier) {
	e.string(id.GetPackageId())
	e.strings(strings.Split(id.GetModuleName(), "."))
	e.strings(strings.Split(id.GetEntityName(), "."))
}

func (e *encoder) create(c *v1.Create, seed []byte, hasSeed bool) {
	e.byte(nodeEncodingVersion)
	e.string(c.GetLfVersion())
	e.byte(0x00)
	if e.present(hasSeed) {
		e.raw(seed)
	}
	e.hexString(c.GetContractId())
	e.string(c.GetPackageName())
	e.identifier(c.GetTemplateId())
	e.value(c.GetArgument())
	e.strings(c.GetSignatories())
	e.strings(c.GetStakeholders())
}

func (e *encoder) fetch(f *v1.Fetch) {
	e.byte(nodeEncodingVersion)
	e.string(f.GetLfVersion())
	e.byte(0x02)
	e.hexString(f.GetContractId())
	e.string(f.GetPackageName())
	e.identifier(f.GetTemplateId())
	e.strings(f.GetSignatories())
	e.strings(f.GetStakeholders())
	if e.present(f.GetInterfaceId() != nil) {
		e.identifier(f.GetInterfaceId())
	}
	e.strings(f.GetActingParties())
}

func (e *encoder) metadata(m *interactive.Metadata) {
	e.byte(metadataEncodingVersion)
	e.strings(m.GetSubmitterInfo().GetActAs())
	e.string(m.GetSubmitterInfo().GetCommandId())
	e.string(m.GetTransactionUuid())
	e.int32(int32(m.GetMediatorGroup()))
	e.string(m.GetSynchronizerId())
	if e.present(m.MinLedgerEffectiveTime != nil) {
		e.int64(int64(m.GetMinLedgerEffectiveTime()))
	}
	if e.present(m.MaxLedgerEffectiveTime != nil) {
		e.int64(int64(m.GetMaxLedgerEffectiveTime()))
	}
	e.int64(int64(m.GetPreparationTime()))
	e.int32(int32(len(m.GetInputContracts())))
	for _, c := range m.GetInputContracts() {
		create := c.GetV1()
		if create == nil {
			e.fail(errors.New("input contract has an unsupported version"))

			return
		}
		e.int64(int64(c.GetCreatedAt()))
		inner := &encoder{}
		inner.create(create, nil, false)
		if inner.err != nil {
			e.fail(fmt.Errorf("input contract %s: %w", create.GetContractId(), inner.err))
		}
		sum := sha256.Sum256(inner.buf)
		e.raw(sum[:])
	}
}

func (e *encoder) value(v *apiv2.Value) {
	switch s := v.GetSum().(type) {
	case *apiv2.Value_Unit:
		e.byte(0x00)
	case *apiv2.Value_Bool:
		e.byte(0x01)
		e.bool(s.Bool)
	case *apiv2.Value_Int64:
		e.byte(0x02)
		e.int64(s.Int64)
	case *apiv2.Value_Numeric:
		e.byte(0x03)
		e.string(s.Numeric)
	case *apiv2.Value_Timestamp:
		e.byte(0x04)
		e.int64(s.Timestamp)
	case *apiv2.Value_Date:
		e.byte(0x05)
		e.int32(s.Date)
	case *apiv2.Value_Party:
		e.byte(0x06)
		e.string(s.Party)
	case *apiv2.Value_Text:
		e.byte(0x07)
		e.string(s.Text)
	case *apiv2.Value_ContractId:
		e.byte(0x08)
		e.hexString(s.ContractId)
	case *apiv2.Value_Optional:
		e.byte(0x09)
		if e.present(s.Optional.GetValue() != nil) {
			e.value(s.Optional.GetValue())
		}
	case *apiv2.Value_List:
		e.byte(0x0a)
		e.int32(int32(len(s.List.GetElements())))
		for _, el := range s.List.GetElements() {
			e.value(el)
		}
	case *apiv2.Value_TextMap:
		e.byte(0x0b)
		e.int32(int32(len(s.TextMap.GetEntries())))
		for _, entry := range s.TextMap.GetEntries() {
			e.string(entry.GetKey())
			e.value(entry.GetValue())
		}
	case *apiv2.Value_Record:
		e.byte(0x0c)
		if e.present(s.Record.GetRecordId() != nil) {
			e.identifier(s.Record.GetRecordId())
		}
		e.int32(int32(len(s.Record.GetFields())))
		for _, f := range s.Record.GetFields() {
			// The spec always encodes the label as set, even when empty.
			e.present(true)
			e.string(f.GetLabel())
			e.value(f.GetValue())
		}
	case *apiv2.Value_Variant:
		e.byte(0x0d)
		if e.present(s.Variant.GetVariantId() != nil) {
			e.identifier(s.Variant.GetVariantId())
		}
		e.string(s.Variant.GetConstructor())
		e.value(s.Variant.GetValue())
	case *apiv2.Value_Enum:
		e.byte(0x0e)
		if e.present(s.Enum.GetEnumId() != nil) {
			e.identifier(s.Enum.GetEnumId())
		}
		e.string(s.Enum.GetConstructor())
	case *apiv2.Value_GenMap:
		e.byte(0x0f)
		e.int32(int32(len(s.GenMap.GetEntries())))
		for _, entry := range s.GenMap.GetEntries() {
			e.value(entry.GetKey())
			e.value(entry.GetValue())
		}
	default:
		e.fail(fmt.Errorf("value %v has no supported type", v))
	}
}
