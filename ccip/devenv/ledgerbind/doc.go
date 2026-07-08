// Package ledgerbind selects CCIP ledger bindings at compile time.
//
// Devenv (default): bindings/generated/latest — current dev DAML layout.
//
// Prod / deployed ledger: go build -tags=prodledger
// Uses bindings/generated/v1_0_0 — template IDs matching Canton TestNet/mainnet today.
package ledgerbind
