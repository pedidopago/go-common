package mariadb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

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
		return []byte("null"), nil
	}
	return json.Marshal(ns.Int64)
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
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
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
		return []byte("null"), nil
	}
	return json.Marshal(ns.Time)
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
		return []byte("null"), nil
	}
	return json.Marshal(ns.Bool)
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
