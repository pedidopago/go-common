package util

import (
	"database/sql/driver"
)

type SecretString string

func (s SecretString) String() string {
	if s == "" {
		return ""
	}
	return "*****"
}

func (s SecretString) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s SecretString) MarshalJSON() ([]byte, error) {
	return []byte("\"" + s.String() + "\""), nil
}

func (s SecretString) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SecretString) Scan(src any) error {
	switch t := src.(type) {
	case string:
		*s = SecretString(t)
	case []byte:
		*s = SecretString(t)
	}
	return nil
}

func NewSecret(s string) SecretString {
	return SecretString(s)
}

// TrueString returns the true value of this string. It is diffenrent than String(),
// because String() returns the safe "*****" value.
func (s SecretString) TrueString() string {
	return string(s)
}
