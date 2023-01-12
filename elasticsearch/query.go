package elasticsearch

type Search struct {
	From  uint64 `json:"from,omitempty"`
	Size  uint64 `json:"size,omitempty"`
	Query Query  `json:"query"`
	Sort  []any  `json:"sort,omitempty"`
}

func (s *Search) Q() *Query {
	return &s.Query
}

func (s *Search) AppendSort(name string, params map[string]any) {
	s.Sort = append(s.Sort, map[string]any{
		name: params,
	})
}

type Query struct {
	Bool     SearchBool      `json:"bool,omitempty"`
	Boosting *SearchBoosting `json:"boosting,omitempty"`
}

type SearchBool struct {
	Must   []map[string]any `json:"must,omitempty"`
	Filter []map[string]any `json:"filter,omitempty"`
	Should []map[string]any `json:"should,omitempty"`
}

type SearchBoosting struct {
	Positive      map[string]any `json:"positive,omitempty"`
	Negative      map[string]any `json:"negative,omitempty"`
	NegativeBoost float64        `json:"negative_boost,omitempty"`
}

type Range map[string]any

func NewRange() Range {
	r := make(map[string]any)
	return Range(r)
}

func RangeGte[T any](s Range, value T) {
	s["gte"] = value
}

func RangeLte[T any](s Range, value T) {
	s["lte"] = value
}

// Term is an exact query
// Match is a fuzzy query

func BoolMustMatch[T any](q *Query, fieldName string, value T) {
	if q.Bool.Must == nil {
		q.Bool.Must = []map[string]any{}
	}
	q.Bool.Must = append(q.Bool.Must, map[string]any{
		"match": map[string]any{
			fieldName: value,
		},
	})
}

func BoolFilterMatch[T any](q *Query, fieldName string, value T) {
	if q.Bool.Filter == nil {
		q.Bool.Filter = []map[string]any{}
	}
	q.Bool.Filter = append(q.Bool.Filter, map[string]any{
		"match": map[string]any{
			fieldName: value,
		},
	})
}

func BoolMustTerm[T any](q *Query, fieldName string, value T) {
	if q.Bool.Must == nil {
		q.Bool.Must = []map[string]any{}
	}
	q.Bool.Must = append(q.Bool.Must, map[string]any{
		"term": map[string]any{
			fieldName: value,
		},
	})
}

func BoolFilterTerm[T any](q *Query, fieldName string, value T) {
	if q.Bool.Filter == nil {
		q.Bool.Filter = []map[string]any{}
	}
	q.Bool.Filter = append(q.Bool.Filter, map[string]any{
		"term": map[string]any{
			fieldName: value,
		},
	})
}

func BoolMustRange(q *Query, fieldName string, rng Range) {
	if q.Bool.Must == nil {
		q.Bool.Must = []map[string]any{}
	}
	rnp := make(map[string]Range)
	rnp[fieldName] = rng
	q.Bool.Must = append(q.Bool.Must, map[string]any{
		"range": rnp,
	})
}

func BoolFilterRange(q *Query, fieldName string, rng Range) {
	if q.Bool.Filter == nil {
		q.Bool.Filter = []map[string]any{}
	}
	rnp := make(map[string]Range)
	rnp[fieldName] = rng
	q.Bool.Filter = append(q.Bool.Filter, map[string]any{
		"range": rnp,
	})
}

//

func BoolShouldMatch[T any](q *Query, fieldName string, value T) {
	if q.Bool.Should == nil {
		q.Bool.Should = []map[string]any{}
	}
	q.Bool.Should = append(q.Bool.Should, map[string]any{
		"match": map[string]any{
			fieldName: value,
		},
	})
}

func BoolShouldTerm[T any](q *Query, fieldName string, value T) {
	if q.Bool.Should == nil {
		q.Bool.Should = []map[string]any{}
	}
	q.Bool.Should = append(q.Bool.Should, map[string]any{
		"term": map[string]any{
			fieldName: value,
		},
	})
}

//

func NewOrderDesc() map[string]any {
	return map[string]any{
		"order": "desc",
	}
}

func NewOrderAsc() map[string]any {
	return map[string]any{
		"order": "asc",
	}
}
