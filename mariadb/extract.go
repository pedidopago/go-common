package mariadb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type extractColumnsOfStructOption struct {
	withBackticks []string
}

type ExtractColumnsOfStructOption func(*extractColumnsOfStructOption)

func WithBackticksColumns(cols ...string) ExtractColumnsOfStructOption {
	return func(o *extractColumnsOfStructOption) {
		if o != nil {
			o.withBackticks = cols
		}
	}
}

func getExtractColumnsOfStructOption(options []ExtractColumnsOfStructOption) extractColumnsOfStructOption {
	opts := extractColumnsOfStructOption{}

	for _, o := range options {
		o(&opts)
	}
	return opts
}

func ExtractColumnsOfStruct(tag string, src any, options ...ExtractColumnsOfStructOption) []string {
	if src == nil {
		return nil
	}
	t := reflect.TypeOf(src)
exOfStruct:
	switch t.Kind() {
	case reflect.Struct:
		return extractColumnsOfStruct("", tag, t, options...)
	case reflect.Pointer:
		t = t.Elem()
		goto exOfStruct
	default:
		return nil
	}
}

func extractColumnsOfStruct(prefix, tag string, src reflect.Type, options ...ExtractColumnsOfStructOption) []string {

	opt := getExtractColumnsOfStructOption(options)

	n := src.NumField()
	outItems := make([]string, 0, n)
	for i := 0; i < n; i++ {
		sf := src.Field(i)
		tname := ""
		tv := strings.Split(sf.Tag.Get(tag), ",")[0]
		if tv == "-" {
			continue
		}
		if tv == "" {
			//TODO: check if is inline (?)
			tname = strings.ToLower(sf.Name)
		} else {
			tname = tv
		}

		for _, withBackticks := range opt.withBackticks {
			if tname == withBackticks {
				tname = fmt.Sprintf("`%s`", tv)
			}
		}

		if sf.Type.Kind() == reflect.Struct {
			if isScannable(sf.Type) {
				outItems = append(outItems, prefix+tname)
			} else {
				xi := extractColumnsOfStruct(prefix+tname+".", tag, sf.Type)
				outItems = append(outItems, xi...)
			}
		} else if sf.Type.Kind() == reflect.Ptr {
			if sf.Type.Elem().Kind() == reflect.Struct {
				if isScannable(sf.Type.Elem()) {
					outItems = append(outItems, prefix+tname)
				} else {
					xi := extractColumnsOfStruct(prefix+tname+".", tag, sf.Type.Elem())
					outItems = append(outItems, xi...)
				}
			}
		} else {
			outItems = append(outItems, prefix+tname)
		}
	}
	return outItems
}

func isScannable(t reflect.Type) bool {
	// check if t implements sql.Scanner
	v := reflect.New(t).Interface()
	if _, ok := v.(sql.Scanner); ok {
		return true
	}
	return t == reflect.TypeOf(time.Time{})
}

func ExtractColumnsAndValues(s interface{}, tag string, ignoreFields ...string) (columns []string, values []interface{}, err error) {
	v := reflect.ValueOf(s)
	v = reflect.Indirect(v)
	t := v.Type()
	if v.Kind() != reflect.Struct {
		return nil, nil, errors.New("type is not a struct")
	}
	if tag == "" {
		return nil, nil, errors.New("missing tag")
	}
	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		field := t.Field(i)
		col := field.Tag.Get(tag)
		if col == "" {
			continue
		}
		if v.Field(i).IsZero() {
			continue
		}
		skip := false
		for _, f := range ignoreFields {
			if col == f {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		val := v.Field(i).Interface()
		columns = append(columns, col)
		values = append(values, val)
	}
	return
}
