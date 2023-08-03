package jsoncolumn

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Text marshals to string when marshaling to the database
type Text[T any] struct {
	Data *T
}

func (c *Text[T]) Scan(src any) error {
	if src == nil {
		c.Data = nil
		return nil
	}

	if c.Data == nil {
		var zv T
		c.Data = &zv
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, c.Data)
	case string:
		return json.Unmarshal([]byte(v), c.Data)
	default:
		return fmt.Errorf("invalid type %T", v)
	}
}

func (c Text[T]) Value() (driver.Value, error) {
	if c.Data == nil {
		return nil, nil
	}
	d, err := json.Marshal(c.Data)
	if err != nil {
		return nil, err
	}
	return string(d), nil
}

func (c *Text[T]) UnmarshalJSON(data []byte) error {
	if c.Data == nil {
		var zv T
		c.Data = &zv
	}
	return json.Unmarshal(data, c.Data)
}

func (c Text[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Data)
}

func (c Text[T]) Deref() T {
	if c.Data == nil {
		var zv T
		return zv
	}
	return *c.Data
}
