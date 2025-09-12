package util

import (
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestSecureID(t *testing.T) {
	allUUIDs := make(map[ulid.ULID]bool)
	var allUUIDsLock sync.Mutex

	setnx := func(id ulid.ULID) (success bool) {
		allUUIDsLock.Lock()
		defer allUUIDsLock.Unlock()
		_, ok := allUUIDs[id]
		if ok {
			return false
		}

		allUUIDs[id] = true
		return true
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go randosUUIDs(t, setnx, 1000, wg)
	}

	wg.Wait()
}

func randosUUIDs(t *testing.T, setnx func(ulid.ULID) bool, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		id := SecureID()
		if !setnx(id) {
			t.Errorf("duplicate id %s", id.String())
		}
	}
}
