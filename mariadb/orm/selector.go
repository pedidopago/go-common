package orm

import (
	"reflect"
	"strings"
)

func ExtractSelectColumnsOfStruct(src any, tag string) []string {
	if src == nil {
		return nil
	}
	t := reflect.TypeOf(src)
exSelOfStruct:
	switch t.Kind() {
	case reflect.Struct:
		return extractSelect(t, tag, "")
	case reflect.Pointer:
		t = t.Elem()
		goto exSelOfStruct
	default:
		panic("src is not a struct or a pointer to a struct: " + t.String())
	}
}

func extractSelect(rtype reflect.Type, tag string, prefix string) []string {
	n := rtype.NumField()
	keys := make([]string, 0, n)
	for i := 0; i < n; i++ {
		sf := rtype.Field(i)
		tvs := strings.Split(sf.Tag.Get(tag), ",")
		if tvs[0] == "-" {
			continue
		}
		if len(tvs) == 1 {
			if tvs[0] == "" {
				keys = append(keys, strings.ToLower(sf.Name))
			} else {
				keys = append(keys, tvs[0])
			}
			continue
		}
		params := tagParams(tvs[1:])
		inline := false
		iprefix := ""
		if _, ok := params["inline"]; ok {
			inline = true
		}
		if _, ok := params["noinsert"]; ok {
			continue
		}
		if v, ok := params["selectprefix"]; ok {
			iprefix = v
		}
		selectOverride := ""
		if v, ok := params["select"]; ok {
			selectOverride = v
			if selectOverride == "@" {
				selectOverride = tvs[0]
			}
		}
		if inline {
			ft := sf.Type
		inlineSelectStart:
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				goto inlineSelectStart
			}
			if ft.Kind() != reflect.Struct {
				// no support for non struct inlines yet!
				continue
			}
			xk := extractSelect(ft, tag, prefix+iprefix)
			keys = append(keys, xk...)
			continue
		}
		if selectOverride != "" {
			keys = append(keys, selectOverride+" AS '"+prefix+iprefix+tvs[0]+"'")
		} else {
			keys = append(keys, prefix+iprefix+tvs[0])
		}
	}
	return keys
}
