package elasticsearch

type Search struct {
	From  uint64 `json:"from"`
	Size  uint64 `json:"size"`
	Query Query  `json:"query"`
	Sort  []any  `json:"sort"`
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
	Bool SearchBool `json:"bool"`
}

type SearchBool struct {
	Must   []map[string]any `json:"must,omitempty"`
	Filter []map[string]any `json:"filter,omitempty"`
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

func BoolMustRange[T any](q *Query, fieldName string, rng Range) {
	if q.Bool.Must == nil {
		q.Bool.Must = []map[string]any{}
	}
	for bi := range q.Bool.Must {
		if rv := q.Bool.Must[bi]["range"]; rv != nil {
			rngg := rv.(map[string]Range)
			rngg[fieldName] = rng
			q.Bool.Must[bi]["range"] = rngg
			return
		}
	}
	rnp := make(map[string]Range)
	rnp[fieldName] = rng
	q.Bool.Must = append(q.Bool.Must, map[string]any{
		"range": rnp,
	})
}

func BoolFilterRange[T any](q *Query, fieldName string, rng Range) {
	if q.Bool.Filter == nil {
		q.Bool.Filter = []map[string]any{}
	}
	for bi := range q.Bool.Filter {
		if rv := q.Bool.Filter[bi]["range"]; rv != nil {
			rngg := rv.(map[string]Range)
			rngg[fieldName] = rng
			q.Bool.Filter[bi]["range"] = rngg
			return
		}
	}
	rnp := make(map[string]Range)
	rnp[fieldName] = rng
	q.Bool.Filter = append(q.Bool.Filter, map[string]any{
		"range": rnp,
	})
}

func BoolMustRanges[T any](q *Query, ranges map[string]Range) {
	if q.Bool.Must == nil {
		q.Bool.Must = []map[string]any{}
	}
	for bi := range q.Bool.Must {
		if rv := q.Bool.Must[bi]["range"]; rv != nil {
			rngg := rv.(map[string]Range)
			for rk, rv := range ranges {
				rngg[rk] = rv
			}
			q.Bool.Must[bi]["range"] = rngg
			return
		}
	}
	q.Bool.Must = append(q.Bool.Must, map[string]any{
		"range": ranges,
	})
}

func BoolFilterRanges[T any](q *Query, ranges map[string]Range) {
	if q.Bool.Filter == nil {
		q.Bool.Filter = []map[string]any{}
	}
	for bi := range q.Bool.Filter {
		if rv := q.Bool.Filter[bi]["range"]; rv != nil {
			rngg := rv.(map[string]Range)
			for rk, rv := range ranges {
				rngg[rk] = rv
			}
			q.Bool.Filter[bi]["range"] = rngg
			return
		}
	}
	q.Bool.Filter = append(q.Bool.Filter, map[string]any{
		"range": ranges,
	})
}
