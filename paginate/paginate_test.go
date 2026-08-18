package paginate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPagination_PreservesExistingBehavior characterizes what the function
// already did before PPS-8997 touched it: the package had no tests at all, so
// these cases are the record of the contract callers depend on today.
func TestNewPagination_PreservesExistingBehavior(t *testing.T) {
	testCases := []struct {
		name         string
		page         int64
		limit        int64
		totalItems   int64
		expectNil    bool
		lastPage     int64
		nearPages    []int64
		itemsPerPage int64
	}{
		{
			name:      "page and limit both zero returns nil",
			page:      0,
			limit:     0,
			expectNil: true,
		},
		{
			name:         "total fits in a single page",
			page:         1,
			limit:        10,
			totalItems:   5,
			lastPage:     1,
			nearPages:    []int64{1},
			itemsPerPage: 10,
		},
		{
			name:         "total equal to the limit still is a single page",
			page:         1,
			limit:        10,
			totalItems:   10,
			lastPage:     1,
			nearPages:    []int64{1},
			itemsPerPage: 10,
		},
		{
			name:         "remainder adds one page",
			page:         1,
			limit:        10,
			totalItems:   95,
			lastPage:     10,
			nearPages:    []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			itemsPerPage: 10,
		},
		{
			name:         "exact division does not add a page",
			page:         2,
			limit:        20,
			totalItems:   100,
			lastPage:     5,
			nearPages:    []int64{1, 2, 3, 4, 5},
			itemsPerPage: 20,
		},
		{
			name:         "zero limit falls back to ten items per page",
			page:         1,
			limit:        0,
			totalItems:   50,
			lastPage:     5,
			nearPages:    []int64{1, 2, 3, 4, 5},
			itemsPerPage: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPagination(tc.page, tc.limit, tc.totalItems)
			if tc.expectNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tc.totalItems, got.TotalItems)
			assert.Equal(t, tc.itemsPerPage, got.ItemsPerPage)
			assert.Equal(t, tc.lastPage, got.LastPage)
			assert.Equal(t, tc.nearPages, got.NearPages)
		})
	}
}

// TestPPS8997_NegativeLimitDoesNotPanic reproduces the panic reported in
// PPS-8997: a negative limit left NearPages empty and the function then read
// its last element, indexing [-1].
func TestPPS8997_NegativeLimitDoesNotPanic(t *testing.T) {
	testCases := []struct {
		name       string
		page       int64
		limit      int64
		totalItems int64
	}{
		{name: "negative limit with no items", page: 1, limit: -5, totalItems: 0},
		{name: "limit of minus one", page: 1, limit: -1, totalItems: 0},
		{name: "negative page and negative limit", page: -3, limit: -10, totalItems: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				got := NewPagination(tc.page, tc.limit, tc.totalItems)
				require.NotNil(t, got)
				assert.GreaterOrEqual(t, got.LastPage, int64(1))
			})
		})
	}
}

// TestPPS8997_NegativeLimitTerminates reproduces the OOM reported in PPS-8997:
// a negative limit made `pages` negative and the loop condition was `!= 0`, so
// the decrement never reached zero and append ran until the pod died. Before
// the fix this test does not finish -- run it with `go test -timeout` and a
// `ulimit -v` guard to see it red without taking the machine down.
func TestPPS8997_NegativeLimitTerminates(t *testing.T) {
	testCases := []struct {
		name         string
		limit        int64
		totalItems   int64
		wantLastPage int64
	}{
		{name: "negative limit smaller than total", limit: -5, totalItems: 100, wantLastPage: 10},
		{name: "negative limit with remainder", limit: -5, totalItems: 102, wantLastPage: 11},
		{name: "negative limit with large total", limit: -1, totalItems: 1000, wantLastPage: 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPagination(1, tc.limit, tc.totalItems)
			require.NotNil(t, got)
			assert.Equal(t, int64(defaultItemsPerPage), got.ItemsPerPage)
			assert.Equal(t, tc.wantLastPage, got.LastPage)
			assert.Len(t, got.NearPages, int(tc.wantLastPage))
		})
	}
}

// TestPPS8997_ClampsPageBelowOne records a deliberate behavior change: a page
// below one used to travel into the response as-is.
func TestPPS8997_ClampsPageBelowOne(t *testing.T) {
	got := NewPagination(-3, 10, 100)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.Page)

	unchanged := NewPagination(3, 10, 100)
	require.NotNil(t, unchanged)
	assert.Equal(t, int64(3), unchanged.Page, "a valid page must not be rewritten")
}

// TestPPS8997_DefaultLimitIsConsistent records a deliberate behavior change:
// the single-page early return used to echo the raw limit, so limit=0 reported
// ItemsPerPage=0 on a short list and 10 on a long one, for the same input.
func TestPPS8997_DefaultLimitIsConsistent(t *testing.T) {
	short := NewPagination(1, 0, 5)
	require.NotNil(t, short)
	assert.Equal(t, int64(defaultItemsPerPage), short.ItemsPerPage)

	long := NewPagination(1, 0, 50)
	require.NotNil(t, long)
	assert.Equal(t, int64(defaultItemsPerPage), long.ItemsPerPage)
	assert.Equal(t, short.ItemsPerPage, long.ItemsPerPage, "the default must not depend on the total")
}

// TestPPS8997_CapsNearPages covers the memory ceiling: limit=1 over a large
// table used to allocate one entry per page (20M rows ~= 152MB) and serialize
// all of it. LastPage stays exact so callers keep the real page count.
func TestPPS8997_CapsNearPages(t *testing.T) {
	t.Run("caps the slice but keeps LastPage exact", func(t *testing.T) {
		got := NewPagination(1, 1, 5_000_000)
		require.NotNil(t, got)
		assert.Len(t, got.NearPages, maxNearPages)
		assert.Equal(t, int64(5_000_000), got.LastPage)
		assert.Equal(t, int64(1), got.NearPages[0])
		assert.Equal(t, int64(maxNearPages), got.NearPages[len(got.NearPages)-1])
	})

	t.Run("exactly at the cap nothing is truncated", func(t *testing.T) {
		got := NewPagination(1, 1, maxNearPages)
		require.NotNil(t, got)
		assert.Len(t, got.NearPages, maxNearPages)
		assert.Equal(t, int64(maxNearPages), got.LastPage)
	})
}
