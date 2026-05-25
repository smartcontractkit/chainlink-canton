package mcms

// MCMS datastore qualifiers for Canton triple-owner governance.
// Each qualifier maps to a distinct MCMS contract instance (separate instanceId@owner).
const (
	QualifierCCIPOwner = "ccipOwner"
	QualifierCCVOwner  = "ccvOwner"
	// QualifierRMNOwner is the MCMS instance for RMNRemote deploy and curse operations.
	QualifierRMNOwner = "rmnOwner"
)
