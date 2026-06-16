# prod_testnet scripts

Canton prod_testnet (`9268731218649498074`) ↔ Sepolia (`16015286601757825753`).

Same flow as `scripts/staging/`, but with prod addresses and **hosted EDS** at `https://eds.testnet.ccip.chain.link` (no local EDS required).

Run from repo root:

```sh
cd /Users/sish/Desktop/chainlink-canton-fcr
```

Always prefer the env wrapper (clears stale JWT/static shell overrides):

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/<script> ...
```

## Setup

1. Create `scripts/prod_testnet/.env` locally (gitignored) with Canton auth, EVM private keys, party IDs, and contract addresses. See script `-help` flags and `prodtestnetenv` env var names for required keys.
2. EDS is pre-deployed — default `PROD_TESTNET_CANTON_EDS_URL=https://eds.testnet.ccip.chain.link`. Override only if you need a local EDS for debugging.
3. Ensure your party has a **PerPartyRouter** + **CCIPSender** on Canton (or let the send script create the router). Update `PROD_TESTNET_CANTON_TO_EVM_SENDER_INSTANCE_ID` if your sender instance id differs from `prod-ccipsender`.
4. **Fund CCIP fees:** Canton→Sepolia send pays fees in **Amulet** (DSO admin, priced on FeeQuoter). Prod has no DevNet tap — your party needs an unlocked Amulet holding. Check balance:
   ```sh
   bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/list_prod_testnet_holdings
   ```
   If empty, transfer a small amount of Amulet to `PROD_TESTNET_CANTON_PARTY_ID` from a funded validator/wallet on prod testnet, then retry send.

## Send (ccipSend)

**Sepolia → Canton** (message only — no token unless `-with-token` or `-token-amount`):

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/send_prod_testnet_evm_to_canton \
  -data 'hello from sepolia to canton'
```

`PROD_TESTNET_EVM_TO_CANTON_TOKEN_AMOUNT` in `.env` does **not** enable token transfer by itself.

**Sepolia → Canton with TEST token** (auto `drip()` on BurnMintERC20WithDrip, mints LINK on Canton):

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/send_prod_testnet_evm_to_canton \
  -data 'hello from sepolia to canton with test' \
  -with-token
```

`-with-token` reads `PROD_TESTNET_EVM_TO_CANTON_TOKEN_AMOUNT` from `.env` (default `1e18` = 1 TEST). Or pass an explicit amount: `-token-amount 1000000000000000000`. Disable drip with `-drip=false`.

Note: `PROD_TESTNET_EVM_TO_CANTON_TOKEN_AMOUNT` in `.env` is **ignored** unless you pass `-with-token` or `-token-amount`.

**Canton → Sepolia** (uses hosted EDS for send disclosures):

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/send_prod_testnet_canton_to_evm \
  -data 'hello from canton to sepolia'
```

**Canton → Sepolia with token** (LINK on Canton → TEST on Sepolia):

Your party must hold **LINK** (`link-token` admin `ccipOwner::1220…`) plus **Amulet** for fees. Hosted EDS token-pool API is not configured for prod yet — `.env` sets `PROD_TESTNET_CANTON_SKIP_TOKEN_POOL_EDS=true` to build pool disclosures from the ledger instead.

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/send_prod_testnet_canton_to_evm_token \
  -data 'hello from canton to sepolia with link'
```

Each send prints a `messageID` — keep it for manual execute.

## Manual execute

**EVM → Canton** (indexer + hosted EDS + Canton ledger):

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/manual_execute_prod_testnet_evm_to_canton \
  -message-id 0x<message-id>
```

**Canton → EVM** (indexer + Sepolia OffRamp execute):

```sh
bash scripts/prod_testnet/with_canton_env.sh ./scripts/prod_testnet/manual_execute_prod_testnet_canton_to_evm \
  -message-id 0x<message-id> \
  -tx-gas-limit 8000000 \
  -gas-limit-override 5000000
```

## Scripts

| Script | Direction |
|--------|-----------|
| `send_prod_testnet_evm_to_canton` | Sepolia → Canton send (message or TEST token) |
| `send_prod_testnet_canton_to_evm` | Canton → Sepolia send (message only) |
| `send_prod_testnet_canton_to_evm_token` | Canton → Sepolia send with LINK token |
| `manual_execute_prod_testnet_evm_to_canton` | Sepolia → Canton execute |
| `manual_execute_prod_testnet_canton_to_evm` | Canton → Sepolia execute |
| `list_prod_testnet_holdings` | List token holdings for fee/token debugging |

## Token lane (LINK ↔ TEST)

| Chain | Asset | Pool / token |
|-------|-------|----------------|
| Canton | LINK | BMTP `0x6ac67a53…` |
| Sepolia | TEST | BMTP `0x5185b41F…`, token `0xeEe6675b…` |

Set `PROD_TESTNET_CCV_ADDRESS_REFS_JSON` to your `cld-signing` copy of `domains/ccv/prod_testnet/datastore/address_refs.json` so ApplyChainUpdates can resolve rate limiters if needed.
