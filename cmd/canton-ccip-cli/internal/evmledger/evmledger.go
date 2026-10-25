// Package evmledger provides EVM transaction signing via a Ledger hardware
// wallet, using go-ethereum's usbwallet implementation.
package evmledger

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ParseDerivationPath parses a --ledger flag value into an EVM derivation path.
// It accepts either a full path like m/44'/60'/0'/0/0 or a depth like 42 in
// which case it will replace the hardened account component, e.g.
// m/44'/60'/42'/0/0.
func ParseDerivationPath(pathOrIndex string) (accounts.DerivationPath, error) {
	if index, err := strconv.ParseUint(pathOrIndex, 10, 32); err == nil {
		path := make(accounts.DerivationPath, len(accounts.DefaultBaseDerivationPath))
		copy(path, accounts.DefaultBaseDerivationPath)
		path[2] = 0x80000000 + uint32(index)

		return path, nil
	}

	return accounts.ParseDerivationPath(pathOrIndex)
}

// NewTransactor connects to a Ledger device, derives the account at the given
// derivation path and returns transaction options that sign EVM transactions
// with the Ledger, together with the derived account address and a function to
// release the device connection.
//
// The usbwallet wallet runs a heartbeat goroutine that permanently tears the
// wallet down as soon as a single health-check exchange with the device fails
// (e.g. when the device is briefly busy or unresponsive), which surfaces as
// accounts.ErrWalletClosed on the next signing attempt. The wallet is
// therefore only kept open for the duration of each signing operation and
// reopened afterwards, so a torn-down wallet never affects the next operation.
func NewTransactor(ctx context.Context, pathOrIndex string, chainID *big.Int) (*bind.TransactOpts, common.Address, func(), error) {
	derivationPath, err := ParseDerivationPath(pathOrIndex)
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("invalid ledger derivation path: %w", err)
	}

	fmt.Println("Looking for connected Ledger devices...")
	hub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to create Ledger hub: %w", err)
	}
	wallets := hub.Wallets()
	if len(wallets) == 0 {
		return nil, common.Address{}, nil, fmt.Errorf("no Ledger wallets found. Please connect a Ledger device with the Ethereum app open and try again")
	}
	fmt.Printf("Found %d connected Ledger wallets\n", len(wallets))
	signer := &ledgerSigner{wallet: wallets[0], path: derivationPath, chainID: chainID}

	// Derive the account once to obtain the sender address. Signing operations
	// open the device themselves, so the connection is not kept open.
	account, err := signer.open()
	if err != nil {
		return nil, common.Address{}, nil, err
	}
	if err := signer.wallet.Close(); err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to close Ledger wallet: %w", err)
	}
	fmt.Printf("Using Ledger account %s at derivation path %s\n", account.Address.Hex(), derivationPath.String())

	auth := &bind.TransactOpts{
		From: account.Address,
		Signer: func(_ common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return signer.signTx(tx)
		},
		Context: ctx,
	}

	return auth, account.Address, func() {}, nil
}

// ledgerSigner signs transactions with a single account on a Ledger device.
type ledgerSigner struct {
	wallet  accounts.Wallet // Pre-discovered Ledger device, reopened per operation
	path    accounts.DerivationPath
	chainID *big.Int
}

// open opens the wallet and derives the account at the signer's derivation path.
func (s *ledgerSigner) open() (accounts.Account, error) {
	if err := s.wallet.Open(""); err != nil {
		return accounts.Account{}, fmt.Errorf("failed to open Ledger wallet: %w", err)
	}
	account, err := s.wallet.Derive(s.path, true)
	if err != nil {
		_ = s.wallet.Close()
		return accounts.Account{}, fmt.Errorf("failed to derive address at derivation path %s: %w", s.path.String(), err)
	}

	return account, nil
}

// signTx opens the wallet, signs the transaction and closes the wallet again.
func (s *ledgerSigner) signTx(tx *types.Transaction) (*types.Transaction, error) {
	account, err := s.open()
	if err != nil {
		return nil, err
	}
	signed, err := s.wallet.SignTx(account, tx, s.chainID)
	if errors.Is(err, accounts.ErrWalletClosed) {
		// The heartbeat health-check may have torn the wallet down between the
		// derive and the sign. No confirmation prompt was shown, so retry once
		// with a freshly opened wallet.
		if account, openErr := s.open(); openErr == nil {
			signed, err = s.wallet.SignTx(account, tx, s.chainID)
		}
	}
	closeErr := s.wallet.Close()
	if err != nil {
		return nil, err
	}
	if closeErr != nil {
		return nil, fmt.Errorf("failed to close Ledger wallet: %w", closeErr)
	}

	return signed, nil
}
