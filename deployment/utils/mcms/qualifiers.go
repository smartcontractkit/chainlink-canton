package mcms

// MCMS datastore qualifiers for Canton dual-owner governance.
// Each qualifier maps to a distinct MCMS contract instance (separate instanceId@owner)
// with the same signer configuration.
const (
	QualifierCCIPOwner = "ccipOwner"
	QualifierCCVOwner  = "ccvOwner"
	// QualifierRMNOwner is the MCMS instance for RMNRemote deploy and curse operations.
	// The MCMS owner party is the same ccipOwner decentralized party; only the instance and config differ.
	QualifierRMNOwner = "rmnOwner"
)
