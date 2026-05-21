package mcms

// MCMS datastore qualifiers for Canton dual-owner governance.
// Each qualifier maps to a distinct MCMS contract instance (separate instanceId@owner)
// with the same signer configuration.
const (
	QualifierCCIPOwner = "ccipOwner"
	QualifierCCVOwner  = "ccvOwner"
)
