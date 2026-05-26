package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

func validMessage() oapiCommon.Message {
	return oapiCommon.Message{
		DestinationChainSelector: "12345",
		Payload:                  "0xdeadbeef",
		Receiver:                 "0xabcdef1234567890",
		FeeToken: oapiCommon.InstrumentId{
			Id:    "token-id",
			Admin: "token-admin",
		},
		Executor: struct {
			Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiCommon.MessageExecutorType `json:"type"`
		}{
			Type: oapiCommon.Empty,
		},
	}
}

func TestValidateMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		modify  func(msg *oapiCommon.Message)
		wantErr string
	}{
		{
			name:   "valid message with empty executor",
			modify: func(msg *oapiCommon.Message) {},
		},
		{
			name: "valid message with noExecutor",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.NoExecutor
			},
		},
		{
			name: "valid message with empty executor and address specified (ignored)",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.Empty
				addr := oapiCommon.RawOrHashedAddress{}
				_ = addr.FromRawInstanceAddress("onramp@owner")
				msg.Executor.Address = &addr
			},
		},
		{
			name: "valid message with noExecutor and address specified (ignored)",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.NoExecutor
				addr := oapiCommon.RawOrHashedAddress{}
				_ = addr.FromRawInstanceAddress("onramp@owner")
				msg.Executor.Address = &addr
			},
		},
		{
			name: "valid message with withAddress executor",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.WithAddress
				addr := oapiCommon.RawOrHashedAddress{}
				_ = addr.FromRawInstanceAddress("onramp@owner")
				msg.Executor.Address = &addr
			},
		},
		{
			name: "valid message with token transfer",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "100.5",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "admin-party",
					},
				}
			},
		},
		{
			name: "valid message with payload without 0x prefix",
			modify: func(msg *oapiCommon.Message) {
				msg.Payload = "deadbeef"
			},
		},
		// Destination chain selector errors
		{
			name: "empty destination chain selector",
			modify: func(msg *oapiCommon.Message) {
				msg.DestinationChainSelector = ""
			},
			wantErr: "invalid destination chain selector",
		},
		{
			name: "non-numeric destination chain selector",
			modify: func(msg *oapiCommon.Message) {
				msg.DestinationChainSelector = "abc"
			},
			wantErr: "invalid destination chain selector",
		},
		{
			name: "zero destination chain selector",
			modify: func(msg *oapiCommon.Message) {
				msg.DestinationChainSelector = "0"
			},
			wantErr: "invalid destination chain selector",
		},
		{
			name: "negative destination chain selector",
			modify: func(msg *oapiCommon.Message) {
				msg.DestinationChainSelector = "-1"
			},
			wantErr: "invalid destination chain selector",
		},
		// Payload errors
		{
			name: "invalid payload hex",
			modify: func(msg *oapiCommon.Message) {
				msg.Payload = "0xZZZZ"
			},
			wantErr: "invalid payload",
		},
		{
			name: "invalid payload hex without prefix",
			modify: func(msg *oapiCommon.Message) {
				msg.Payload = "not-hex"
			},
			wantErr: "invalid payload",
		},
		// Receiver errors
		{
			name: "invalid receiver hex",
			modify: func(msg *oapiCommon.Message) {
				msg.Receiver = "0xGGGG"
			},
			wantErr: "invalid receiver",
		},
		// Fee token errors
		{
			name: "empty fee token id",
			modify: func(msg *oapiCommon.Message) {
				msg.FeeToken.Id = ""
			},
			wantErr: "fee token id is required",
		},
		{
			name: "empty fee token admin",
			modify: func(msg *oapiCommon.Message) {
				msg.FeeToken.Admin = ""
			},
			wantErr: "fee token admin is required",
		},
		// Executor errors
		{
			name: "invalid executor type",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = "invalid"
			},
			wantErr: "invalid executor type",
		},
		{
			name: "withAddress executor without address",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.WithAddress
				msg.Executor.Address = nil
			},
			wantErr: "executor address is required",
		},
		{
			name: "withAddress executor with empty address",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.WithAddress
				addr := oapiCommon.RawOrHashedAddress{}
				_ = addr.FromRawInstanceAddress("")
				msg.Executor.Address = &addr
			},
			wantErr: "invalid executor address",
		},
		{
			name: "withAddress executor with invalid address",
			modify: func(msg *oapiCommon.Message) {
				msg.Executor.Type = oapiCommon.WithAddress
				addr := oapiCommon.RawOrHashedAddress{}
				_ = addr.FromRawInstanceAddress("invalid")
				msg.Executor.Address = &addr
			},
			wantErr: "invalid executor address",
		},
		// Token transfer errors
		{
			name: "token transfer with empty token id",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "100",
					Token: oapiCommon.InstrumentId{
						Id:    "",
						Admin: "admin",
					},
				}
			},
			wantErr: "token transfer token id is required",
		},
		{
			name: "token transfer with empty token admin",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "100",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "",
					},
				}
			},
			wantErr: "token transfer token admin is required",
		},
		{
			name: "token transfer with empty amount",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "admin",
					},
				}
			},
			wantErr: "token transfer amount is required",
		},
		{
			name: "token transfer with zero amount",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "0",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "admin",
					},
				}
			},
			wantErr: "token transfer amount must be greater than zero",
		},
		{
			name: "token transfer with zero decimal amount",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "0.0",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "admin",
					},
				}
			},
			wantErr: "token transfer amount must be greater than zero",
		},
		{
			name: "token transfer with negative amount",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "-1.5",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "admin",
					},
				}
			},
			wantErr: "token transfer amount must be greater than zero",
		},
		{
			name: "token transfer with invalid amount",
			modify: func(msg *oapiCommon.Message) {
				msg.TokenTransfer = &oapiCommon.TokenTransfer{
					Amount: "not-a-number",
					Token: oapiCommon.InstrumentId{
						Id:    "link",
						Admin: "admin",
					},
				}
			},
			wantErr: "invalid token transfer amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			msg := validMessage()
			tt.modify(&msg)

			err := ValidateMessage(msg)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
