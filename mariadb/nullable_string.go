package mariadb

import (
	"database/sql"
	"database/sql/driver"
)

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
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + ns.String + `"`), nil
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	ns.String = string(data)
	return nil
}
