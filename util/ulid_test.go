package util

import (
	"sync"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
)

func TestSecureID(t *testing.T) {
	uuidchan := make(chan ulid.ULID, 10000)
	allUUIDs := make(map[ulid.ULID]bool)
	var allUUIDsLock sync.RWMutex

	addid := func(id ulid.ULID) {
		uuidchan <- id
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go randosUUIDs(t, addid, 100, wg)
	}

	go func() {
		for id := range uuidchan {
			allUUIDsLock.RLock()
			_, ok := allUUIDs[id]
			allUUIDsLock.RUnlock()
			if ok {
				t.Errorf("duplicate UUID: %s", id)
			} else {
				allUUIDsLock.Lock()
				allUUIDs[id] = true
				allUUIDsLock.Unlock()
			}
		}
	}()

	wg.Wait()

	time.Sleep(time.Millisecond * 100)

	close(uuidchan)

	if len(allUUIDs) != 500000 {
		t.Errorf("expected 500000 UUIDs, got %d", len(allUUIDs))
	}
}

func randosUUIDs(t *testing.T, addid func(ulid.ULID), n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		id := SecureID()
		addid(id)
	}
}
