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
