// Package elasticsearch provides a set of helper functions to the elasticsearch API
//
// Deprecated: this package is deprecated.
package elasticsearch

// Deprecated: this package is deprecated.
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

func (s *Search) AppendSortNoParams(name string) {
	s.Sort = append(s.Sort, name)
}

// Deprecated: this package is deprecated.
type Query struct {
	Bool     SearchBool      `json:"bool,omitempty"`
	Boosting *SearchBoosting `json:"boosting,omitempty"`
}

// Deprecated: this package is deprecated.
type SearchBool struct {
	Must    []map[string]any `json:"must,omitempty"`
	MustNot []map[string]any `json:"must_not,omitempty"`
	Filter  []map[string]any `json:"filter,omitempty"`
	Should  []map[string]any `json:"should,omitempty"`
}

// Deprecated: this package is deprecated.
type SearchBoosting struct {
	Positive      map[string]any `json:"positive,omitempty"`
	Negative      map[string]any `json:"negative,omitempty"`
	NegativeBoost float64        `json:"negative_boost,omitempty"`
}

// Deprecated: this package is deprecated.
type SearchWildcard map[string]any

// Deprecated: this package is deprecated.
type Range map[string]any

// Deprecated: this package is deprecated.
func NewRange() Range {
	r := make(map[string]any)
	return Range(r)
}

// Deprecated: this package is deprecated.
func RangeGte[T any](s Range, value T) {
	s["gte"] = value
}

// Deprecated: this package is deprecated.
func RangeLte[T any](s Range, value T) {
	s["lte"] = value
}

// Term is an exact query
// Match is a fuzzy query

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
func BoolFilterWildcardDefaults(q *Query, fieldName string, value string) {
	BoolFilterWildcard(q, fieldName, value, 1.0, "constant_score")
}

// Deprecated: this package is deprecated.
func BoolShouldWildcardDefaults(q *Query, fieldName string, value string) {
	BoolShouldWildcard(q, fieldName, value, 1.0, "constant_score")
}

// Deprecated: this package is deprecated.
type WildcardStruct struct {
	Value   string  `json:"value"`
	Boost   float64 `json:"boost,omitempty"`
	Rewrite string  `json:"rewrite,omitempty"`
}

// Deprecated: this package is deprecated.
func BoolFilterWildcard(q *Query, fieldName string, value string, boost float64, rewrite string) {
	if q.Bool.Filter == nil {
		q.Bool.Filter = []map[string]any{}
	}
	q.Bool.Filter = append(q.Bool.Filter, map[string]any{
		"wildcard": map[string]any{
			fieldName: WildcardStruct{
				Value:   value,
				Boost:   boost,
				Rewrite: rewrite,
			},
		},
	})
}

// Deprecated: this package is deprecated.
func BoolShouldWildcard(q *Query, fieldName string, value string, boost float64, rewrite string) {
	if q.Bool.Should == nil {
		q.Bool.Should = []map[string]any{}
	}
	q.Bool.Should = append(q.Bool.Should, map[string]any{
		"wildcard": map[string]any{
			fieldName: WildcardStruct{
				Value:   value,
				Boost:   boost,
				Rewrite: rewrite,
			},
		},
	})
}

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
func BoolMust(q *Query, fieldName string, value any) {
	if q.Bool.Must == nil {
		q.Bool.Must = []map[string]any{}
	}
	q.Bool.Must = append(q.Bool.Must, map[string]any{
		fieldName: value,
	})
}

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
func BoolMustNotTerm[T any](q *Query, fieldName string, value T) {
	if q.Bool.MustNot == nil {
		q.Bool.MustNot = []map[string]any{}
	}
	q.Bool.MustNot = append(q.Bool.MustNot, map[string]any{
		"term": map[string]any{
			fieldName: value,
		},
	})
}

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
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

// Deprecated: this package is deprecated.
func NewOrderDesc() map[string]any {
	return map[string]any{
		"order": "desc",
	}
}

// Deprecated: this package is deprecated.
func NewOrderAsc() map[string]any {
	return map[string]any{
		"order": "asc",
	}
}
