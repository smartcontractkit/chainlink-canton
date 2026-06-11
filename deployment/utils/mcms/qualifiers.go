package mcms

// MCMS datastore qualifiers for Canton triple-owner governance.
// Each qualifier maps to a distinct MCMS contract instance (separate instanceId@owner).
//
// QualifierCCIPOwner uses "CLLCCIP" (not the Canton owner party name "ccipOwner") so Canton
// and EVM chains share the same MCMS lookup key in address_refs.json. Cross-family merged
// proposals (e.g. Canton ↔ Sepolia lane configure) require one mcms.qualifier across all
// chains in the batch; EVM prod uses CLLCCIP. The on-ledger signatory remains the
// ccipOwner:: party — only the datastore label changes. deploy_mcms already defaults empty
// qualifier to CLLCCIP; staging_testnet_canton uses the same pattern.
const (
	QualifierCCIPOwner = "CLLCCIP"
	QualifierCCVOwner  = "ccvOwner"
	// QualifierRMNOwner uses "RMNMCMS" (not the Canton owner party name "rmnOwner") so Canton
	// and EVM chains share the same MCMS lookup key in address_refs.json. Future cross-family RMN
	// proposals (e.g. Curse) require one mcms.qualifier across all chains in the batch; EVM prod
	// uses RMNMCMS. The on-ledger signatory remains the rmnOwner:: party — only the datastore
	// label changes.
	QualifierRMNOwner = "RMNMCMS"
)
