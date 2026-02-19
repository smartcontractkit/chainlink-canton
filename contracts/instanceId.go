package contracts

import (
	"fmt"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"math/rand/v2"
	"strings"
)

const (
	// InstanceIDPrefixHintMaxLength defines the maximum length for the prefix hint in an InstanceID.
	InstanceIDPrefixHintMaxLength = 50
	// InstanceIDRandomPartLength defines the length of the random prefix component in an InstanceID.
	InstanceIDRandomPartLength = 5
)

// InstanceID represents a unique identifier for a contract instance.
type InstanceID string

// NewInstanceID creates a new InstanceID using the provided prefixHint.
// The format of the InstanceID is "<prefixHint>-<randomPart>".
// The prefixHint is converted to lowercase and must only contain alphanumeric characters, hyphens, or underscores.
// The randomPart is a randomly generated string of lowercase letters of length InstanceIDRandomPartLength.
func NewInstanceID(prefixHint string) (InstanceID, error) {
	if prefixHint == "" {
		return "", fmt.Errorf("instance ID prefixHint cannot be empty")
	}
	if len(prefixHint) > InstanceIDPrefixHintMaxLength {
		// This is an arbitrary restriction to prevent overly long instance IDs.
		return "", fmt.Errorf("instance ID prefix hint cannot be longer than %d characters", InstanceIDPrefixHintMaxLength)
	}
	// Check that prefixHint contains only valid characters (alphanumeric and hyphens/underscores)
	for i, char := range prefixHint {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '-' && char != '_' {
			return "", fmt.Errorf("instance ID prefix hint contains invalid character '%c' at position %d", char, i)
		}
	}

	randomPrefix := generateRandomComponent()
	instanceId := InstanceID(fmt.Sprintf("%s-%s", strings.ToLower(prefixHint), randomPrefix))

	return instanceId, nil
}

// MustNewInstanceID is like NewInstanceID but panics if an error occurs.
func MustNewInstanceID(prefixHint string) InstanceID {
	instanceID, err := NewInstanceID(prefixHint)
	if err != nil {
		panic(fmt.Sprintf("failed to create InstanceID: %v", err))
	}

	return instanceID
}

// RawInstanceAddress returns the RawInstanceAddress for this InstanceID and the given party.
func (i InstanceID) RawInstanceAddress(party types.PARTY) RawInstanceAddress {
	return RawInstanceAddress(fmt.Sprintf("%s@%s", i, party))
}

// String returns the string representation of the InstanceID.
func (i InstanceID) String() string {
	return string(i)
}

func generateRandomComponent() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	s := make([]rune, InstanceIDRandomPartLength)
	for i := range s {
		//nolint:gosec // Not used for cryptographic purposes
		s[i] = rune(charset[rand.IntN(len(charset))])
	}

	return string(s)
}
