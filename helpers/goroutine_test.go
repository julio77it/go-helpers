package helpers_test

import (
	"sync"
	"testing"

	"github.com/julio77it/go-helpers/helpers"
)

func TestGetGID(t *testing.T) {
	numgo := 10

	ch := make(chan uint64, numgo)

	var wg sync.WaitGroup

	wg.Add(numgo)

	for i := 0; i < numgo; i++ {
		go func() {
			ch <- helpers.GetGID()
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	collected := make(map[uint64]uint64)

	for v := range ch {
		collected[v] = v
	}
	if len(collected) != numgo {
		t.Errorf("GetGID distinct goroutines. expected %d, got %d", numgo, len(collected))
	}
}
