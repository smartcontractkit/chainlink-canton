package parse

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/converters"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

func RawInstanceAddress(a chainlinkapi.RawInstanceAddress) (contracts.RawInstanceAddress, error) {
	address, err := contracts.RawInstanceAddressFromString(string(a.Unpack))
	if err != nil {
		return contracts.RawInstanceAddress(""), fmt.Errorf("failed to parse raw instance address: %w", err)
	}

	return address, nil
}

func RawInstanceAddressList(addrs []chainlinkapi.RawInstanceAddress) ([]contracts.RawInstanceAddress, error) {
	addresses := make([]contracts.RawInstanceAddress, len(addrs))
	for i, addr := range addrs {
		parsedAddr, err := RawInstanceAddress(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse raw instance address at index %d: %w", i, err)
		}
		addresses[i] = parsedAddr
	}

	return addresses, nil
}

// ValidateMessage validates the fields of a CCIP message.
func ValidateMessage(msg oapiCommon.Message) error {
	// Destination chain selector must be a valid non-zero uint64
	if chainSelector, err := strconv.ParseUint(msg.DestinationChainSelector, 10, 64); err != nil || chainSelector == 0 {
		return fmt.Errorf("invalid destination chain selector: %w", err)
	}

	// Payload must be a valid hex string (with optional 0x prefix)
	if _, err := hex.DecodeString(strings.TrimPrefix(msg.Payload, "0x")); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	// Receiver must be a valid hex string (with optional 0x prefix)
	if _, err := hex.DecodeString(strings.TrimPrefix(msg.Receiver, "0x")); err != nil {
		return fmt.Errorf("invalid receiver: %w", err)
	}

	// Fee token must have a valid instrument identifier
	if msg.FeeToken.Id == "" {
		return fmt.Errorf("fee token id is required")
	}
	if msg.FeeToken.Admin == "" {
		return fmt.Errorf("fee token admin is required")
	}

	// Executor type determines whether an address is required
	switch msg.Executor.Type {
	case oapiCommon.Empty, oapiCommon.NoExecutor:
		// No address expected, ignored if specified
	case oapiCommon.WithAddress:
		if msg.Executor.Address == nil {
			return fmt.Errorf("executor address is required when type is %q", oapiCommon.WithAddress)
		}
		if _, err := converters.ResolveRawOrHashedAddress(*msg.Executor.Address); err != nil {
			return fmt.Errorf("invalid executor address: %w", err)
		}
	default:
		return fmt.Errorf("invalid executor type: %q", msg.Executor.Type)
	}

	// Token transfer fields (optional, validated only when present)
	if msg.TokenTransfer != nil {
		if msg.TokenTransfer.Token.Id == "" {
			return fmt.Errorf("token transfer token id is required")
		}
		if msg.TokenTransfer.Token.Admin == "" {
			return fmt.Errorf("token transfer token admin is required")
		}
		if msg.TokenTransfer.Amount == "" {
			return fmt.Errorf("token transfer amount is required")
		}
		amount, err := strconv.ParseFloat(msg.TokenTransfer.Amount, 64)
		if err != nil {
			return fmt.Errorf("invalid token transfer amount: %w", err)
		}
		if amount <= 0 {
			return fmt.Errorf("token transfer amount must be greater than zero")
		}
	}

	return nil
}
