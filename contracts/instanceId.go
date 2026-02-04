package contracts

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"

	"golang.org/x/crypto/sha3"
)

const (
	// InstanceIDPrefixHintMaxLength defines the maximum length for the prefix hint in an InstanceID.
	InstanceIDPrefixHintMaxLength = 50
	// InstanceIDRandomPartLength defines the length of the random prefix component in an InstanceID.
	InstanceIDRandomPartLength = 5
)

// InstanceID represents a unique identifier for a contract instance.
type InstanceID string

// NewInstanceID creates a new InstanceID using the provided prefixHint and party.
// The format of the InstanceID is "<prefixHint>-<randomPart>@<party>".
// The prefixHint is converted to lowercase and must only contain alphanumeric characters, hyphens, or underscores.
// The randomPart is a randomly generated string of lowercase letters of length InstanceIDRandomPartLength.
// The party is included as-is and must not be empty.
func NewInstanceID(prefixHint string, party string) (InstanceID, error) {
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
	if party == "" {
		return "", errors.New("party cannot be empty")
	}

	randomPrefix := generateRandomComponent()
	instanceId := InstanceID(fmt.Sprintf("%s-%s@%s", strings.ToLower(prefixHint), randomPrefix, party))

	return instanceId, nil
}

// MustNewInstanceID is like NewInstanceID but panics if an error occurs.
func MustNewInstanceID(prefixHint string, party string) InstanceID {
	instanceID, err := NewInstanceID(prefixHint, party)
	if err != nil {
		panic(fmt.Sprintf("failed to create InstanceID: %v", err))
	}

	return instanceID
}

// InstanceAddress computes the InstanceAddress for the InstanceID by hashing it using Keccak256.
func (i InstanceID) InstanceAddress() InstanceAddress {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(i))

	return InstanceAddress(h.Sum(nil))
}

// String returns the string representation of the InstanceID.
func (i InstanceID) String() string {
	return string(i)
}

// Valid checks if the InstanceID is in the correct format.
// A valid InstanceID has the format "<prefix>@<party>" where both prefix and party are non-empty.
func (i InstanceID) Valid() bool {
	parts := strings.Split(string(i), "@")
	if len(parts) != 2 {
		return false
	}

	return len(parts[0]) > 0 && len(parts[1]) > 0
}

// Party extracts the party component from the InstanceID.
func (i InstanceID) Party() (string, error) {
	parts := strings.Split(string(i), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid instance ID format")
	}

	return parts[1], nil
}

// Prefix extracts the prefix component from the InstanceID.
func (i InstanceID) Prefix() (string, error) {
	parts := strings.Split(string(i), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid instance ID format")
	}

	return parts[0], nil
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
