package mariadb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"
	"unicode/utf8"
)

var (
	jsonNull  = []byte("null")
	jsonTrue  = []byte("true")
	jsonFalse = []byte("false")
)

const jsonHex = "0123456789abcdef"

func appendJSONString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	start := 0
	for i := 0; i < len(s); {
		c := s[i]
		if c >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r != utf8.RuneError && r != ' ' && r != ' ' {
				i += size
				continue
			}
			dst = append(dst, s[start:i]...)
			dst = append(dst, '\\', 'u')
			dst = append(dst, jsonHex[(r>>12)&0xf], jsonHex[(r>>8)&0xf], jsonHex[(r>>4)&0xf], jsonHex[r&0xf])
			i += size
			start = i
			continue
		}
		if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
			i++
			continue
		}
		dst = append(dst, s[start:i]...)
		switch c {
		case '"':
			dst = append(dst, '\\', '"')
		case '\\':
			dst = append(dst, '\\', '\\')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', jsonHex[c>>4], jsonHex[c&0xf])
		}
		i++
		start = i
	}
	dst = append(dst, s[start:]...)
	dst = append(dst, '"')
	return dst
}

// OpenAPISchemaObject is implemented by types that can describe themselves as an OpenAPI schema.
type OpenAPISchemaObject interface {
	SetType(v string)
	SetFormat(v string)
	SetDescription(v string)
}

// NullInt64 represents an int64 that may be NULL in SQL and null in JSON.
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullInt64) Scan(value any) error {
	if value == nil {
		ns.Int64, ns.Valid = 0, false
		return nil
	}
	ns2 := &sql.NullInt64{}
	if err := ns2.Scan(value); err != nil {
		return err
	}
	ns.Valid = true
	ns.Int64 = ns2.Int64
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullInt64) Value() (driver.Value, error) {
	return sql.NullInt64{
		Int64: ns.Int64,
		Valid: ns.Valid,
	}.Value()
}

// MarshalJSON implements the json.Marshaler interface. Null values are encoded as JSON null.
func (ns NullInt64) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	return strconv.AppendInt(make([]byte, 0, 20), ns.Int64, 10), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. JSON null sets Valid to false.
func (ns *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.Int64)
}

// Set sets the value and marks it as valid (non-NULL).
func (ns *NullInt64) Set(v int64) { ns.Int64 = v; ns.Valid = true }

// Clear resets the value to zero and marks it as NULL.
func (ns *NullInt64) Clear() { ns.Int64 = 0; ns.Valid = false }

// IsZero returns true if the value is NULL (not valid).
func (ns NullInt64) IsZero() bool { return !ns.Valid }

// HydrateSchemaObject populates an OpenAPI schema describing this type.
func (ns NullInt64) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("integer")
	schema.SetFormat("int64")
	schema.SetDescription("nullable int64")
}

// NullString represents a string that may be NULL in SQL and null in JSON.
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value any) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns2 := &sql.NullString{}
	if err := ns2.Scan(value); err != nil {
		return err
	}
	ns.Valid = true
	ns.String = ns2.String
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	return sql.NullString{
		String: ns.String,
		Valid:  ns.Valid,
	}.Value()
}

// MarshalJSON implements the json.Marshaler interface. Null values are encoded as JSON null.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	return appendJSONString(make([]byte, 0, len(ns.String)+2), ns.String), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. JSON null sets Valid to false.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}

// Set sets the value and marks it as valid (non-NULL).
func (ns *NullString) Set(v string) { ns.String = v; ns.Valid = true }

// Clear resets the value to empty and marks it as NULL.
func (ns *NullString) Clear() { ns.String = ""; ns.Valid = false }

// IsZero returns true if the value is NULL (not valid).
func (ns NullString) IsZero() bool { return !ns.Valid }

// HydrateSchemaObject populates an OpenAPI schema describing this type.
func (ns NullString) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("string")
	schema.SetDescription("nullable string")
}

// String returns a valid NullString for non-empty strings, or an invalid one for empty strings.
func String(s string) NullString {
	return NullString{
		String: s,
		Valid:  s != "",
	}
}

// NullTime represents a time.Time that may be NULL in SQL and null in JSON.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTime) Scan(value any) error {
	if value == nil {
		ns.Time, ns.Valid = time.Time{}, false
		return nil
	}
	ns2 := &sql.NullTime{}
	if err := ns2.Scan(value); err != nil {
		return err
	}
	ns.Valid = true
	ns.Time = ns2.Time
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	return sql.NullTime{
		Time:  ns.Time,
		Valid: ns.Valid,
	}.Value()
}

// MarshalJSON implements the json.Marshaler interface. Null values are encoded as JSON null.
func (ns NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	return ns.Time.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface. JSON null sets Valid to false.
func (ns *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.Time)
}

// Set sets the value and marks it as valid (non-NULL).
func (ns *NullTime) Set(v time.Time) { ns.Time = v; ns.Valid = true }

// Clear resets the value to zero time and marks it as NULL.
func (ns *NullTime) Clear() { ns.Time = time.Time{}; ns.Valid = false }

// IsZero returns true if the value is NULL (not valid).
func (ns NullTime) IsZero() bool { return !ns.Valid }

// HydrateSchemaObject populates an OpenAPI schema describing this type.
func (ns NullTime) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("string")
	schema.SetDescription("nullable RFC3339 date-time")
}

// ToTimePtr returns a pointer to the time value, or nil if NULL.
func (ns NullTime) ToTimePtr() *time.Time {
	if !ns.Valid {
		return nil
	}
	return &ns.Time
}

// Time returns a valid NullTime for non-zero times, or an invalid one for zero times.
func Time(t time.Time) NullTime {
	return NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// NullBool represents a bool that may be NULL in SQL and null in JSON.
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBool) Scan(value any) error {
	if value == nil {
		ns.Bool, ns.Valid = false, false
		return nil
	}
	ns2 := &sql.NullBool{}
	if err := ns2.Scan(value); err != nil {
		return err
	}
	ns.Valid = true
	ns.Bool = ns2.Bool
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullBool) Value() (driver.Value, error) {
	return sql.NullBool{
		Bool:  ns.Bool,
		Valid: ns.Valid,
	}.Value()
}

// MarshalJSON implements the json.Marshaler interface. Null values are encoded as JSON null.
func (ns NullBool) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	if ns.Bool {
		return jsonTrue, nil
	}
	return jsonFalse, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. JSON null sets Valid to false.
func (ns *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.Bool)
}

// Set sets the value and marks it as valid (non-NULL).
func (ns *NullBool) Set(v bool) { ns.Bool = v; ns.Valid = true }

// Clear resets the value to false and marks it as NULL.
func (ns *NullBool) Clear() { ns.Bool = false; ns.Valid = false }

// IsZero returns true if the value is NULL (not valid).
func (ns NullBool) IsZero() bool { return !ns.Valid }

// HydrateSchemaObject populates an OpenAPI schema describing this type.
func (ns NullBool) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("boolean")
	schema.SetDescription("nullable boolean")
}

// FromString returns a NullString. If zeroIsNil is true, an empty string produces a NULL value.
func FromString(v string, zeroIsNil bool) NullString {
	if zeroIsNil && v == "" {
		return NullString{}
	}
	return NullString{
		String: v,
		Valid:  true,
	}
}

// FromTime returns a NullTime. If zeroIsNil is true, a zero time produces a NULL value.
func FromTime(t time.Time, zeroIsNil bool) NullTime {
	if zeroIsNil && t.IsZero() {
		return NullTime{}
	}
	return NullTime{
		Time:  t,
		Valid: true,
	}
}

// FromBool returns a NullBool. If zeroIsNil is true, false produces a NULL value.
func FromBool(b bool, zeroIsNil bool) NullBool {
	if zeroIsNil && !b {
		return NullBool{}
	}
	return NullBool{
		Bool:  b,
		Valid: true,
	}
}

// FromInt64 returns a NullInt64. If zeroIsNil is true, zero produces a NULL value.
func FromInt64(i int64, zeroIsNil bool) NullInt64 {
	if zeroIsNil && i == 0 {
		return NullInt64{}
	}
	return NullInt64{
		Int64: i,
		Valid: true,
	}
}
