package devenv

import (
	"bytes"
	"encoding/binary"
)

var GenericExtraArgsV3Tag = []byte{0xa6, 0x9d, 0xd4, 0xaa}

type GenericExtraArgsV3 struct {
	GasLimit           uint32
	BlockConfirmations uint16
	CCVs               [][]byte
	CCVArgs            [][]byte
	Executor           []byte
	ExecutorArgs       []byte
	TokenReceiver      []byte
	TokenArgs          []byte
}

func EncodeGenericExtraArgsV3(args *GenericExtraArgsV3) ([]byte, error) {
	var buf bytes.Buffer

	// Write tag
	buf.Write(GenericExtraArgsV3Tag)

	// Write gasLimit (uint32, big-endian)
	_ = binary.Write(&buf, binary.BigEndian, args.GasLimit)

	// Write blockConfirmations (uint16, big-endian)
	_ = binary.Write(&buf, binary.BigEndian, args.BlockConfirmations)

	// Write CCVs length
	_ = buf.WriteByte(uint8(len(args.CCVs))) //nolint:gosec

	// Write each CCV
	for i, v := range args.CCVs {
		// Write CCV address length (uint8)
		buf.WriteByte(uint8(len(v))) //nolint:gosec
		// Write CCV address itself
		buf.Write(v)

		// Write CCV args length (uint16, big-endian)
		_ = binary.Write(&buf, binary.BigEndian, uint16(len(args.CCVArgs[i]))) //nolint:gosec
		// Write CCV arg itself
		buf.Write(args.CCVArgs[i])
	}

	// Write Executor
	buf.WriteByte(uint8(len(args.Executor))) //nolint:gosec
	buf.Write(args.Executor)

	// Write Executor args
	_ = binary.Write(&buf, binary.BigEndian, uint16(len(args.ExecutorArgs))) //nolint:gosec
	buf.Write(args.ExecutorArgs)

	// Write Token receiver
	buf.WriteByte(uint8(len(args.TokenReceiver))) //nolint:gosec
	buf.Write(args.TokenReceiver)

	// Write Token Args
	_ = binary.Write(&buf, binary.BigEndian, uint16(len(args.TokenArgs))) //nolint:gosec
	buf.Write(args.TokenArgs)

	return buf.Bytes(), nil
}
