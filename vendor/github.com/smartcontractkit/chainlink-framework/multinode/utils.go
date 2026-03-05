package multinode

import (
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/jpillora/backoff"
)

func RandomID() ID {
	// #nosec G404
	id := rand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

func NewIDFromInt(id int64) ID {
	return big.NewInt(id)
}

// NewRedialBackoff is a standard backoff to use for redialling or reconnecting to
// unreachable network endpoints
func NewRedialBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}
