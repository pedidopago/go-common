package mariadb

import (
	"database/sql"
	"reflect"
	"strings"
	"time"
)

func ExtractColumnsOfStruct(tag string, src any) []string {
	if src == nil {
		return nil
	}
	t := reflect.TypeOf(src)
exOfStruct:
	switch t.Kind() {
	case reflect.Struct:
		return extractColumnsOfStruct("", tag, t)
	case reflect.Pointer:
		t = t.Elem()
		goto exOfStruct
	default:
		return nil
	}
}

func extractColumnsOfStruct(prefix, tag string, src reflect.Type) []string {
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
