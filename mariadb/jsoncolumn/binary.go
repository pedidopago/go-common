package jsoncolumn

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Binary marshals to []byte when marshaling to database
type Binary[T any] struct {
	Data *T
}

func (c *Binary[T]) Scan(src any) error {
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

func (c Binary[T]) Value() (driver.Value, error) {
	if c.Data == nil {
		return nil, nil
	}
	return json.Marshal(c.Data)
}

func (c *Binary[T]) UnmarshalJSON(data []byte) error {
	if c.Data == nil {
		var zv T
		c.Data = &zv
	}
	return json.Unmarshal(data, c.Data)
}

func (c Binary[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Data)
}
