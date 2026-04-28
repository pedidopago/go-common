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

type OpenAPISchemaObject interface {
	SetType(v string)
	SetFormat(v string)
	SetDescription(v string)
}

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

func (ns NullInt64) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	return strconv.AppendInt(make([]byte, 0, 20), ns.Int64, 10), nil
}

func (ns *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.Int64)
}

func (ns NullInt64) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("integer")
	schema.SetFormat("int64")
	schema.SetDescription("nullable int64")
}

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

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	return appendJSONString(make([]byte, 0, len(ns.String)+2), ns.String), nil
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}

func (ns NullString) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("string")
	schema.SetDescription("nullable string")
}

func String(s string) NullString {
	return NullString{
		String: s,
		Valid:  s != "",
	}
}

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

func (ns NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	return ns.Time.MarshalJSON()
}

func (ns *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.Time)
}

func (ns NullTime) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("string")
	schema.SetDescription("nullable RFC3339 date-time")
}

func (ns NullTime) ToTimePtr() *time.Time {
	if !ns.Valid {
		return nil
	}
	return &ns.Time
}

func Time(t time.Time) NullTime {
	return NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

//

type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Int64 is not NULL
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

func (ns NullBool) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return jsonNull, nil
	}
	if ns.Bool {
		return jsonTrue, nil
	}
	return jsonFalse, nil
}

func (ns *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.Bool)
}

func (ns NullBool) HydrateSchemaObject(schema OpenAPISchemaObject) {
	schema.SetType("boolean")
	schema.SetDescription("nullable boolean")
}
