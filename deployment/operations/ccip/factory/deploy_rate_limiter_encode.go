package factory

import (
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
)

// encodeDeployRateLimiterParamsWire encodes DeployRateLimiter MCMS params for Factory.daml.
// Factory.decodeDeployRateLimiterParams expects direction/mode as single-byte tags (0/1),
// but go-daml HexCodec encodes enums as length-prefixed constructor names. All other fields
// use the standard HexCodec wire format.
func encodeDeployRateLimiterParamsWire(params factorybindings.DeployRateLimiterParams) (string, error) {
	hexCodec := codec.NewHexCodec()

	appendEncoded := func(value any) ([]byte, error) {
		part, err := hexCodec.Marshal(value)
		if err != nil {
			return nil, err
		}

		return []byte(part), nil
	}

	var wire []byte

	for _, value := range []any{
		params.InstanceId,
		params.PoolInstanceId,
		params.PoolOwner,
		params.RemoteChainSelector,
	} {
		part, err := appendEncoded(value)
		if err != nil {
			return "", fmt.Errorf("encode DeployRateLimiter field: %w", err)
		}
		wire = append(wire, part...)
	}

	direction, err := rateLimitDirectionWireByte(params.Direction)
	if err != nil {
		return "", err
	}
	wire = append(wire, direction)

	mode, err := rateLimitModeWireByte(params.Mode)
	if err != nil {
		return "", err
	}
	wire = append(wire, mode)

	for _, value := range []any{params.IsEnabled, params.Capacity, params.Rate} {
		part, err := appendEncoded(value)
		if err != nil {
			return "", fmt.Errorf("encode DeployRateLimiter field: %w", err)
		}
		wire = append(wire, part...)
	}

	return hex.EncodeToString(wire), nil
}

func rateLimitDirectionWireByte(direction core.RateLimitDirection) (byte, error) {
	switch direction {
	case core.RateLimitDirectionRateLimitDirection_Outbound:
		return 0, nil
	case core.RateLimitDirectionRateLimitDirection_Inbound:
		return 1, nil
	default:
		return 0, fmt.Errorf("unknown rate limit direction: %q", direction)
	}
}

func rateLimitModeWireByte(mode core.RateLimitMode) (byte, error) {
	switch mode {
	case core.RateLimitModeRateLimitMode_DefaultFinality:
		return 0, nil
	case core.RateLimitModeRateLimitMode_CustomFinality:
		return 1, nil
	default:
		return 0, fmt.Errorf("unknown rate limit mode: %q", mode)
	}
}

func encodeDeployRateLimiter(args factorybindings.DeployRateLimiter) (*bind.EncodedChoice, error) {
	c := args.Contract
	opData, err := encodeDeployRateLimiterParamsWire(factorybindings.DeployRateLimiterParams{
		InstanceId:          c.InstanceId,
		PoolInstanceId:      c.PoolInstanceId,
		PoolOwner:           c.PoolOwner,
		RemoteChainSelector: c.RemoteChainSelector,
		Direction:           c.Direction,
		Mode:                c.Mode,
		IsEnabled:           c.IsEnabled,
		Capacity:            c.Capacity,
		Rate:                c.Rate,
	})
	if err != nil {
		return nil, err
	}

	return &bind.EncodedChoice{
		Choice:        "DeployRateLimiter",
		OperationData: opData,
	}, nil
}
