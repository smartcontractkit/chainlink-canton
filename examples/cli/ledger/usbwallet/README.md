# usbwallet (Canton)

This package is a Canton-specific Ledger USB wallet backend used by the example CLI and tests.
It is heavily adapted from go-ethereum's `accounts/usbwallet` code, but speaks the Canton Ledger
app protocol and signs Ed25519 payloads for Canton workflows.

## What it provides

- Ledger device discovery over HID (`NewCantonLedgerHub`, `Wallets`)
- Device/session helpers (`Open`, `Close`, `Status`, `GetVersion`, `GetAppName`)
- Public key retrieval and Canton fingerprinting (`GetPublicKey`, `Fingerprint`)
- Signing flows for:
  - raw hash / multihash (`SignHash`)
  - prepared Canton interactive submissions (`SignPreparedTransaction`)
  - onboarding topology transactions, optionally with attestation challenge (`SignTopologyTransactions`)

## Canton-specific behavior

- Uses Canton app APDUs (`INS_GET_VERSION`, `INS_GET_PUBLIC_KEY`, `INS_SIGN_TX`, ...)
- Uses Ed25519 keys/signatures (not ECDSA)
- Streams large payloads in chunked APDUs (max payload per APDU: 255 bytes)
- For prepared submissions, splits the transaction into device-expected protobuf messages
  (transaction header, ordered nodes, metadata, input contracts)
- Preserves unknown `driver_metadata` fields in input contracts so host/device hash computation stays aligned

## Prerequisites

- A Ledger device with the Canton app open
- Blind signing enabled on the Canton app
- A valid derivation path (for Canton typically `m/44'/6767'/0'/0'/0'`)

## Minimal usage

```go
derivationPath, _ := accounts.ParseDerivationPath("m/44'/6767'/0'/0'/0'")

hub, err := usbwallet.NewCantonLedgerHub()
if err != nil {
	panic(err)
}

wallets := hub.Wallets()
if len(wallets) == 0 {
	panic("no Ledger wallet found")
}
wallet := wallets[0]
if err := wallet.Open(""); err != nil {
	panic(err)
}
defer wallet.Close()

pubKey, err := wallet.GetPublicKey(derivationPath, false)
if err != nil {
	panic(err)
}
fingerprint, _ := usbwallet.Fingerprint(pubKey)
_ = fingerprint
```
