package binder

import (
	"encoding"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// CustomBinder is a variation of echo.DefaultBinder
type CustomBinder struct {
	echoDefaultBinder echo.DefaultBinder
	SplitCharacter    string
	StopCharacter     string
}

// ensure our CustomBinder implements echo.Binder
var _ echo.Binder = &CustomBinder{}

// ensure CustomBinder implements EchoContextBinder
var _ EchoContextBinder = &CustomBinder{}

type EchoContextBindUnmarshaler interface {
	UnmarshalParam(c echo.Context, param string) error
}

func (b *CustomBinder) Bind(i any, c echo.Context) (err error) {
	if err := b.BindPathParams(i, c); err != nil {
		return err
	}
	if err = b.BindQueryParams(i, c); err != nil {
		return err
	}
	return b.BindBody(i, c)
}

// BindQueryParams provides a customized implementation of this method
func (b *CustomBinder) BindQueryParams(i any, c echo.Context) error {
	if err := b.bindData(c, i, c.QueryParams(), "query", "mquery"); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}
	return nil
}

func (b *CustomBinder) BindPathParams(i any, c echo.Context) error {
	return b.echoDefaultBinder.BindPathParams(c, i)
}

func (b *CustomBinder) BindBody(i any, c echo.Context) error {
	return b.echoDefaultBinder.BindBody(c, i)
}

// bindData is a customized version that allows a field to have more than one alias
func (b *CustomBinder) bindData(c echo.Context, destination interface{}, data map[string][]string, tag string, customTags ...string) error {
	if destination == nil || len(data) == 0 {
		return nil
	}
	typ := reflect.TypeOf(destination).Elem()
	val := reflect.ValueOf(destination).Elem()

	// Map
	if typ.Kind() == reflect.Map {
		for k, v := range data {
			val.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
		}
		return nil
	}

	// !struct
	if typ.Kind() != reflect.Struct {
		if tag == "param" || tag == "query" || tag == "header" {
			// incompatible type, data is probably to be found in the body
			return nil
		}
		return errors.New("binding element must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if typeField.Anonymous {
			if structField.Kind() == reflect.Ptr {
				structField = structField.Elem()
			}
		}
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)
		if typeField.Anonymous && structField.Kind() == reflect.Struct && inputFieldName != "" {
			// if anonymous struct with query/param/form tags, report an error
			return errors.New("query/param/form tags are not allowed with anonymous struct field")
		}

		if inputFieldName == "" {
			// If tag is nil, we inspect if the field is a not BindUnmarshaler struct and try to bind data into it (might contains fields with tags).
			// structs that implement BindUnmarshaler are binded only when they have explicit tag
			if _, ok := structField.Addr().Interface().(echo.BindUnmarshaler); !ok && structFieldKind == reflect.Struct {
				if err := b.bindData(c, structField.Addr().Interface(), data, tag, customTags...); err != nil {
					return err
				}
			}
			// does not have explicit tag and is not an ordinary struct - so move to next field
			continue
		}

		inputValue, exists := data[inputFieldName]
		if !exists && len(customTags) != 0 {
			inputFieldNameMain := inputFieldName

		outer:
			for _, t := range customTags {
				tv := typeField.Tag.Get(t)
				var stopAt = strings.Index(tv, b.StopCharacter)
				if stopAt == -1 {
					stopAt = len(tv)
				}
				for _, inputFieldName = range strings.Split(tv[:stopAt], b.SplitCharacter) {
					inputValue, exists = data[inputFieldName]
					if exists {
						break outer
					}
				}
			}
			if !exists {
				inputFieldName = inputFieldNameMain
			}
		}

		if !exists {
			// Go json.Unmarshal supports case insensitive binding.  However the
			// url params are bound case sensitive which is inconsistent.  To
			// fix this we must check all of the map values in a
			// case-insensitive search.
			for k, v := range data {
				if strings.EqualFold(k, inputFieldName) {
					inputValue = v
					exists = true
					break
				}
			}
		}

		if !exists {
			continue
		}

		// Call this first, in case we're dealing with an alias to an array type
		if ok, err := unmarshalField(c, typeField.Type.Kind(), inputValue, structField); ok {
			if err != nil {
				return err
			}
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithProperType(c, sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if err := setWithProperType(c, typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err

		}
	}
	return nil
}

func unmarshalField(c echo.Context, valueKind reflect.Kind, val []string, field reflect.Value) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return unmarshalFieldPtr(c, val, field)
	default:
		return unmarshalFieldNonPtr(c, val, field)
	}
}

func unmarshalFieldNonPtr(c echo.Context, values []string, field reflect.Value) (bool, error) {
	fieldIValue := field.Addr().Interface()
	var callbackFn func(value string) error
	if unmarshaler, ok := fieldIValue.(EchoContextBindUnmarshaler); ok {
		callbackFn = func(value string) error {
			return unmarshaler.UnmarshalParam(c, value)
		}
	}
	if unmarshaler, ok := fieldIValue.(echo.BindUnmarshaler); ok {
		callbackFn = func(value string) error {
			return unmarshaler.UnmarshalParam(value)
		}
	}
	if unmarshaler, ok := fieldIValue.(encoding.TextUnmarshaler); ok {
		callbackFn = func(value string) error {
			return unmarshaler.UnmarshalText([]byte(value))
		}
	}

	if callbackFn == nil {
		return false, nil
	}

	for _, value := range values {
		if err := callbackFn(value); err != nil {
			return true, nil
		}
	}

	return true, nil
}

func unmarshalFieldPtr(c echo.Context, value []string, field reflect.Value) (bool, error) {
	if field.IsNil() {
		// Initialize the pointer to a nil value
		field.Set(reflect.New(field.Type().Elem()))
	}
	return unmarshalFieldNonPtr(c, value, field.Elem())
}

func setWithProperType(c echo.Context, valueKind reflect.Kind, val string, structField reflect.Value) error {
	// But also call it here, in case we're dealing with an array of BindUnmarshalers
	if ok, err := unmarshalField(c, valueKind, []string{val}, structField); ok {
		return err
	}

	switch valueKind {
	case reflect.Ptr:
		return setWithProperType(c, structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
