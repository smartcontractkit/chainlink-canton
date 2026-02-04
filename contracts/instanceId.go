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
	InstanceIDRandomPartLength    = 5
)

type InstanceID string

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

func MustNewInstanceID(prefixHint string, party string) InstanceID {
	instanceID, err := NewInstanceID(prefixHint, party)
	if err != nil {
		panic(fmt.Sprintf("failed to create InstanceID: %v", err))
	}

	return instanceID
}

func (i InstanceID) InstanceAddress() InstanceAddress {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(i))

	return InstanceAddress(h.Sum(nil))
}

func (i InstanceID) String() string {
	return string(i)
}

func (i InstanceID) Valid() bool {
	parts := strings.Split(string(i), "@")
	if len(parts) != 2 {
		return false
	}

	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func (i InstanceID) Party() (string, error) {
	parts := strings.Split(string(i), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid instance ID format")
	}

	return parts[1], nil
}

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
