package changesets

type CantonCSDeps[CFG any] struct {
	ChainSelector uint64

	// Which participant in the environment to use (0-indexed, i.e. defaults to the first participant)
	Participant int

	// TODO Should be part of the env, not exposed by CLDF/CTF yet
	UserName string
	Party    string

	Config CFG
}
