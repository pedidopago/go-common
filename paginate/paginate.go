package paginate

const (
	defaultItemsPerPage = 10
	// maxNearPages caps how many page numbers one response carries. A caller
	// asking for limit=1 over a multi-million row table would otherwise get one
	// entry per page -- 20M rows is ~152MB of int64 before serialization -- so
	// the ceiling is a memory guard, not a display choice. LastPage stays exact,
	// so the real page count survives the cap.
	maxNearPages = 10000
)

type Pagination struct {
	TotalItems   int64   `json:"total_items"`
	ItemsPerPage int64   `json:"items_per_page"`
	Page         int64   `json:"page"`
	LastPage     int64   `json:"last_page"`
	NearPages    []int64 `json:"near_pages"`
}

func NewPagination(page, limit, totalItems int64) *Pagination {
	// AIDEV-NOTE: this nil is published contract -- callers serialize the result
	// straight to JSON -- so it must stay ahead of the clamps below (PPS-8997).
	if page == 0 && limit == 0 {
		return nil
	}
	// AIDEV-NOTE: page and limit arrive from the query string, and MySQL coerces
	// a negative LIMIT/OFFSET to unsigned instead of rejecting it, so a hostile
	// value reaches this function intact (PPS-8997).
	if limit < 1 {
		limit = defaultItemsPerPage
	}
	if page < 1 {
		page = 1
	}
	pagination := &Pagination{
		TotalItems:   totalItems,
		ItemsPerPage: limit,
		Page:         page,
	}
	if totalItems <= limit {
		pagination.LastPage = 1
		pagination.NearPages = []int64{1}
		return pagination
	}
	pages := totalItems / limit
	if totalItems%limit > 0 {
		pages++
	}
	// AIDEV-NOTE: LastPage comes from the page count, never from the last element
	// of NearPages -- that indexing is what panicked on an empty slice, and the
	// count has to outlive the maxNearPages cap (PPS-8997).
	pagination.LastPage = pages
	near := min(pages, maxNearPages)
	pagination.NearPages = make([]int64, 0, near)
	for index := int64(1); index <= near; index++ {
		pagination.NearPages = append(pagination.NearPages, index)
	}
	return pagination
}
