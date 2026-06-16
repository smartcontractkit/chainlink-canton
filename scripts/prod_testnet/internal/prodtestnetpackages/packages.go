// Package prodtestnetpackages maps CCIP ledger template IDs to prod testnet bundled DAR package names.
package prodtestnetpackages

import (
	"fmt"
	"strings"
	"sync"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/scripts/prod_testnet/internal/prodtestnetenv"
)

var (
	loadOnce sync.Once
	pkgs     packageNames
)

type packageNames struct {
	Core                 string
	Runtime              string
	Receiver             string
	Sender               string
	Executor             string
	BurnMintTokenPool    string
	LockReleaseTokenPool string
}

type templateBinding interface {
	GetTemplateID() string
}

// Init loads prod ledger package names from env (call after prodtestnetenv.LoadDefault).
func Init() {
	loadOnce.Do(func() {
		pkgs = packageNames{
			Core:                 envPackage("ccip-core", "PROD_TESTNET_CORE_PACKAGE"),
			Runtime:              envPackage("ccip-runtime", "PROD_TESTNET_RUNTIME_PACKAGE", "PROD_TESTNET_PER_PARTY_ROUTER_PACKAGE"),
			Receiver:             envPackage("ccip-receiver", "PROD_TESTNET_CCIP_RECEIVER_PACKAGE"),
			Sender:               envPackage("ccip-sender", "PROD_TESTNET_CCIP_SENDER_PACKAGE"),
			Executor:             envPackage("ccip-executor", "PROD_TESTNET_EXECUTOR_PACKAGE"),
			BurnMintTokenPool:    envPackage("ccip-burn-mint-token-pool", "PROD_TESTNET_BURN_MINT_TOKEN_POOL_PACKAGE"),
			LockReleaseTokenPool: envPackage("ccip-lock-release-token-pool", "PROD_TESTNET_LOCK_RELEASE_TOKEN_POOL_PACKAGE", "PROD_TESTNET_TOKEN_POOL_PACKAGE"),
		}
	})
}

func envPackage(defaultName string, keys ...string) string {
	return strings.TrimPrefix(strings.TrimSpace(prodtestnetenv.String(defaultName, keys...)), "#")
}

// ledgerPackageNameRef is the PackageId form for Canton upgradable package names (#name syntax).
func ledgerPackageNameRef(name string) string {
	name = strings.TrimPrefix(strings.TrimSpace(name), "#")
	if name == "" {
		return ""
	}
	return "#" + name
}

func templateID(ledgerPackage string, template templateBinding) string {
	base := contracts.TemplateIDFromBinding(template).String()
	pkg := strings.TrimPrefix(strings.TrimSpace(ledgerPackage), "#")
	if pkg == "" {
		return base
	}
	// ReplacePackageIdWithNameInTemplateID yields "ccip-runtime:Module:Entity"; ledger filters need "#ccip-runtime:…".
	replaced := contracts.ReplacePackageIdWithNameInTemplateID(base, pkg)
	parts := strings.SplitN(replaced, ":", 3)
	if len(parts) != 3 {
		return replaced
	}
	return fmt.Sprintf("%s:%s:%s", ledgerPackageNameRef(parts[0]), parts[1], parts[2])
}

func ledgerTemplate(ledgerPackage string, template templateBinding) *apiv2.Identifier {
	tid, err := contracts.TemplateIDFromString(templateID(ledgerPackage, template))
	if err != nil {
		panic(err)
	}
	return tid.ToLedgerIdentifier()
}

// CoreTemplateID is for templates bundled in ccip-core (GlobalConfig, RateLimiter, FeeQuoter, …).
func CoreTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.Core, template)
}

// RuntimeTemplateID is for ccip-runtime (PerPartyRouter, OnRamp, OffRamp, …).
func RuntimeTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.Runtime, template)
}

func ReceiverTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.Receiver, template)
}

func SenderTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.Sender, template)
}

func ExecutorTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.Executor, template)
}

func BurnMintTokenPoolTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.BurnMintTokenPool, template)
}

func LockReleaseTokenPoolTemplateID(template templateBinding) string {
	Init()
	return templateID(pkgs.LockReleaseTokenPool, template)
}

// TokenPoolTemplateID is an alias for LockReleaseTokenPoolTemplateID (legacy env PROD_TESTNET_TOKEN_POOL_PACKAGE).
func TokenPoolTemplateID(template templateBinding) string {
	return LockReleaseTokenPoolTemplateID(template)
}

func GlobalConfigLedgerTemplate() *apiv2.Identifier {
	Init()
	return ledgerTemplate(pkgs.Core, core.GlobalConfig{})
}

func PerPartyRouterFactoryLedgerTemplate() *apiv2.Identifier {
	Init()
	return ledgerTemplate(pkgs.Runtime, ccipruntime.PerPartyRouterFactory{})
}

func PerPartyRouterLedgerTemplate() *apiv2.Identifier {
	Init()
	return ledgerTemplate(pkgs.Runtime, ccipruntime.PerPartyRouter{})
}

func CCIPReceiverLedgerTemplate() *apiv2.Identifier {
	Init()
	return ledgerTemplate(pkgs.Receiver, receiver.CCIPReceiver{})
}
