package paginate

type Pagination struct {
	TotalItems   int64   `json:"total_items"`
	ItemsPerPage int64   `json:"items_per_page"`
	Page         int64   `json:"page"`
	LastPage     int64   `json:"last_page"`
	NearPages    []int64 `json:"near_pages"`
}

func NewPagination(page, limit, totalItems int64) *Pagination {
	if page == 0 && limit == 0 {
		return nil
	}
	pagination := &Pagination{
		TotalItems:   totalItems,
		ItemsPerPage: limit,
		Page:         page,
	}
	if limit == 0 {
		pagination.ItemsPerPage = 10
	}
	if totalItems <= pagination.ItemsPerPage {
		return &Pagination{
			TotalItems:   totalItems,
			ItemsPerPage: limit,
			Page:         page,
			LastPage:     1,
			NearPages:    []int64{1},
		}
	}
	pages, remainder := pagination.TotalItems/pagination.ItemsPerPage, pagination.TotalItems%pagination.ItemsPerPage
	if remainder > 0 {
		pages++
	}
	var index int64 = 1
	for pages != 0 {
		pagination.NearPages = append(pagination.NearPages, index)
		index++
		pages--
	}
	pagination.LastPage = pagination.NearPages[len(pagination.NearPages)-1]
	return pagination
}
