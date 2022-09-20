package mariadb

import (
	"encoding/json"
	"testing"
)

func TestNullString(t *testing.T) {
	ns1j := `{"a": null, "b": "ok"}`
	ns1s := struct {
		A NullString `json:"a"`
		B NullString `json:"b"`
	}{}
	if err := json.Unmarshal([]byte(ns1j), &ns1s); err != nil {
		t.Fatal(err)
	}
	if ns1s.A.Valid {
		t.Fatal("A should be invalid")
	}
	if !ns1s.B.Valid {
		t.Fatal("B should be valid")
	}
	if ns1s.B.String != "ok" {
		t.Fatal("B should be 'ok'")
	}
}

func TestNullTime(t *testing.T) {
	ns1j := `{"a": null, "b": "2009-11-10T23:00:00Z"}`
	ns1s := struct {
		A NullTime `json:"a"`
		B NullTime `json:"b"`
	}{}
	if err := json.Unmarshal([]byte(ns1j), &ns1s); err != nil {
		t.Fatal(err)
	}
	if ns1s.A.Valid {
		t.Fatal("A should be invalid")
	}
	if !ns1s.B.Valid {
		t.Fatal("B should be valid")
	}
	if ns1s.B.Time.Year() != 2009 {
		t.Fatal("B should be '2009-11-10T23:00:00Z'")
	}
}
