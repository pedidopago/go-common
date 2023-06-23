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

// StringIsULID returns true if the string is a valid ULID
func StringIsULID(s string) bool {
	if len(s) != 26 {
		return false
	}
	_, err := ulid.Parse(s)
	return err == nil
}

func init() {
	secureSource = ulid.Monotonic(srand.Reader, 0)
}
