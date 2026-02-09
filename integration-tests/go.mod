module github.com/smartcontractkit/chainlink-canton/integration-tests

go 1.25.5

replace github.com/smartcontractkit/chainlink-canton => ../

replace github.com/digital-asset/dazl-client/v8 => github.com/noders-team/dazl-client/v8 v8.7.1-2

// Coming from chainlink-deployments-framework
replace github.com/fbsobreira/gotron-sdk => github.com/smartcontractkit/chainlink-tron/relayer/gotron-sdk v0.0.5-0.20251014120029-d73d15cc23f7

replace github.com/smartcontractkit/go-daml => github.com/smartcontractkit/go-daml v0.0.0-20260209193910-eac8a15b0b35

require (
	github.com/digital-asset/dazl-client/v8 v8.8.0
	github.com/ethereum/go-ethereum v1.16.8
	github.com/google/uuid v1.6.0
	github.com/smartcontractkit/chainlink-canton v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/ProjectZKM/Ziren/crates/go-runtime/zkvm_runtime v0.0.0-20251001021608-1fe7b43fc4d6 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/oapi-codegen/runtime v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	gonum.org/v1/gonum v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
