package api

type RequestMeta struct {
	Page       int         `json:"page" schema:"page"`
	SortFields []SortField `json:"sort" schema:"sort"`
	// Requested count of records per page
	// -1 == all records;
	PerPage int           `json:"per_page" schema:"per_page"`
	Skip    int           `json:"-" schema:"-"`
	Filters []FilterField `json:"filters" schema:"filters"`
}

type FilterField struct {
	Key   string
	Value string
}

// Check method sets page to 1 if it less than 1
// and sets skip argument
func (r *RequestMeta) Check() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Page > 1 {
		r.Skip = (r.Page - 1) * r.PerPage
	}
}

type ResponseMeta struct {
	// Current page
	Page int `json:"page"`
	// Count of records on this page
	Count int `json:"count"`
	// Total count of records could be paginated
	Total int `json:"total"`
}

type Sort []SortField

type SortField struct {
	Name       string `json:"name"`
	Descending bool   `json:"descending"`
}
