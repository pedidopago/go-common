package util

import (
	srand "crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var secureSource *ulid.MonotonicEntropy
var idmutex sync.Mutex

// SecureID returns a Universally Unique Lexicographically Sortable Identifier
// obtained via crypto/rand entropy
func SecureID() ulid.ULID {
	idmutex.Lock()
	defer idmutex.Unlock()
	return ulid.MustNew(ulid.Timestamp(time.Now()), secureSource)
}

func init() {
	secureSource = ulid.Monotonic(srand.Reader, 0)
}
