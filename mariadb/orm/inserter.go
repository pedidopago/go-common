package orm

import (
	"reflect"
	"strings"
)

// ExtractInsertColumnsOfStruct extracts the columns of a struct that are
// suitable for insertion. This is useful for structs that have fields that
// are not suitable for insertion (e.g. auto-incrementing primary keys) but
// are required for other operations.
//
// Panics if src is not a struct or a pointer to a struct.
func ExtractInsertColumnsOfStruct(src any, tag string) (keys []string, values []any) {
	if src == nil {
		return nil, nil
	}
	t := reflect.TypeOf(src)
	v := reflect.ValueOf(src)
exOfStruct:
	switch t.Kind() {
	case reflect.Struct:
		return extractInsert(t, v, tag)
	case reflect.Pointer:
		t = t.Elem()
		v = v.Elem()
		goto exOfStruct
	default:
		panic("src is not a struct or a pointer to a struct: " + t.String())
	}
}

func extractInsert(rtype reflect.Type, value reflect.Value, tag string) (keys []string, values []any) {
	n := rtype.NumField()
	keys = make([]string, 0, n)
	values = make([]any, 0, n)
	for i := 0; i < n; i++ {
		sf := rtype.Field(i)
		tvs := strings.Split(sf.Tag.Get(tag), ",")
		if tvs[0] == "-" {
			continue
		}
		if len(tvs) == 1 {
			if tvs[0] == "" {
				tvs[0] = strings.ToLower(sf.Name)
			} else {
				keys = append(keys, tvs[0])
			}
			values = append(values, value.Field(i).Interface())
			continue
		}
		params := tagParams(tvs[1:])
		omitempty := false
		inline := false
		if _, ok := params["omitempty"]; ok {
			omitempty = true
		}
		if _, ok := params["inline"]; ok {
			inline = true
		}
		if _, ok := params["noinsert"]; ok {
			continue
		}
		if value.Field(i).IsZero() && omitempty {
			continue
		}
		if inline {
			fieldValue := value.Field(i)
		inlineInsertStart:
			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				continue
			}
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
				goto inlineInsertStart
			}
			if fieldValue.Kind() != reflect.Struct {
				// no support for non struct inlines yet!
				continue
			}
			xk, xv := extractInsert(sf.Type, value.Field(i), tag)
			keys = append(keys, xk...)
			values = append(values, xv...)
			continue
		}
		keys = append(keys, tvs[0])
		values = append(values, value.Field(i).Interface())
	}
	return keys, values
}

func tagParams(tslc []string) map[string]string {
	output := make(map[string]string)
	for _, ts := range tslc {
		if ts == "" {
			continue
		}
		kv := strings.SplitN(ts, "=", 2)
		if len(kv) == 1 {
			output[kv[0]] = ""
		} else {
			output[kv[0]] = kv[1]
		}
	}
	return output
}
