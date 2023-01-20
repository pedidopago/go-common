package types

import (
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

type String string

func (s String) String() string {
	return string(s)
}

func (String) FromString(v string) (String, error) {
	return String(v), nil
}

func (String) SchemaTypeKind() reflect.Kind { return reflect.String }

type IntType[T constraints.Integer] struct {
	value T
}

func (IntType[T]) FromString(v string) (IntType[T], error) {
	p, err := strconv.ParseInt(v, 10, 64)
	return IntType[T]{value: T(p)}, err
}

func (i IntType[T]) Int() T {
	return i.value
}

func (i IntType[T]) SchemaTypeKind() reflect.Kind { return reflect.TypeOf(i.value).Kind() }

type Boolean bool

func (Boolean) FromString(v string) (Boolean, error) {
	p, err := strconv.ParseBool(v)
	return Boolean(p), err
}

func (b Boolean) Bool() bool {
	return bool(b)
}

func (Boolean) SchemaTypeKind() reflect.Kind { return reflect.Bool }

type FloatType[T constraints.Float] struct {
	value T
}

func (FloatType[T]) FromString(v string) (FloatType[T], error) {
	p, err := strconv.ParseFloat(v, 64)
	return FloatType[T]{value: T(p)}, err
}

func (f FloatType[T]) Float() T {
	return f.value
}

func (f FloatType[T]) SchemaTypeKind() reflect.Kind { return reflect.TypeOf(f.value).Kind() }
