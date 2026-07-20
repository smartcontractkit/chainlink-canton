# ledgertarget

Compile-time adapter for Canton ledger layout differences.

Local devenv uses `bindings/generated/latest` (current DAML). Prod testnet/mainnet still runs `v1_0_0` contracts with different module paths and field names.

Go build tags select the target: default for devenv, `-tags=prodledger` for prod. All tag-split bindings and operations live here so the rest of `devenv` stays tag-free.
