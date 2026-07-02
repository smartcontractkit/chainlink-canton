package client

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsTopologyMappingAlreadyExists reports whether err is Canton's signal that an
// identical topology mapping is already authorized on the synchronizer.
func IsTopologyMappingAlreadyExists(err error) bool {
	for err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
			if strings.Contains(st.Message(), "TOPOLOGY_MAPPING_ALREADY_EXISTS") {
				return true
			}
		}
		if strings.Contains(err.Error(), "TOPOLOGY_MAPPING_ALREADY_EXISTS") {
			return true
		}
		err = errors.Unwrap(err)
	}

	return false
}

// IsKmsPublicKeyConflict reports whether err is Canton's signal that the KMS
// key material is already registered in the vault under a different name.
func IsKmsPublicKeyConflict(err error) bool {
	for err != nil {
		msg := err.Error()
		if strings.Contains(msg, "Existing public key") && strings.Contains(msg, "different than inserted key") {
			return true
		}
		if st, ok := status.FromError(err); ok {
			msg = st.Message()
			if strings.Contains(msg, "Existing public key") && strings.Contains(msg, "different than inserted key") {
				return true
			}
		}
		err = errors.Unwrap(err)
	}

	return false
}
