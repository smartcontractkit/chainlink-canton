package contracts

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
)

const (
	InstanceIDPrefixMaxLength  = 6
	InstanceIDRandomPartLength = 5
)

type InstanceID string

func NewInstanceID(prefixHint string, party string) (InstanceID, error) {
	if prefixHint == "" {
		return "", fmt.Errorf("instance ID prefixHint cannot be empty")
	}
	if len(prefixHint) > InstanceIDPrefixMaxLength {
		return "", fmt.Errorf("instance ID prefixHint cannot be longer than %d characters", InstanceIDPrefixMaxLength)
	}
	if party == "" {
		return "", errors.New("party cannot be empty")
	}

	randomPrefix := generateRandomComponent()
	instanceId := InstanceID(fmt.Sprintf("%s%s-%s@%s", strings.ToLower(prefixHint), strings.Repeat("0", InstanceIDPrefixMaxLength-len(prefixHint)), randomPrefix, party))

	return instanceId, nil
}

func MustNewInstanceID(prefixHint string, party string) InstanceID {
	instanceID, err := NewInstanceID(prefixHint, party)
	if err != nil {
		panic(fmt.Sprintf("failed to create InstanceID: %v", err))
	}

	return instanceID
}

func (i InstanceID) String() string {
	return string(i)
}

func (i InstanceID) Valid() bool {
	parts := strings.Split(string(i), "@")
	if len(parts) != 2 {
		return false
	}

	return len(parts[0]) == InstanceIDPrefixMaxLength+1+InstanceIDRandomPartLength
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
		s[i] = rune(charset[rand.IntN(len(charset))])
	}
	return string(s)
}
