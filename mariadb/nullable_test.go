package mariadb

import (
	"encoding/json"
	"math"
	"testing"
	"time"
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

func TestNullInt64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		val  NullInt64
		want string
	}{
		{"null", NullInt64{}, "null"},
		{"zero", NullInt64{Int64: 0, Valid: true}, "0"},
		{"positive", NullInt64{Int64: 42, Valid: true}, "42"},
		{"negative", NullInt64{Int64: -1, Valid: true}, "-1"},
		{"max", NullInt64{Int64: math.MaxInt64, Valid: true}, "9223372036854775807"},
		{"min", NullInt64{Int64: math.MinInt64, Valid: true}, "-9223372036854775808"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.val.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNullBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		val  NullBool
		want string
	}{
		{"null", NullBool{}, "null"},
		{"true", NullBool{Bool: true, Valid: true}, "true"},
		{"false", NullBool{Bool: false, Valid: true}, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.val.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNullString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		val   NullString
		want  string
	}{
		{"null", NullString{}, "null"},
		{"empty", NullString{String: "", Valid: true}, `""`},
		{"simple", NullString{String: "hello", Valid: true}, `"hello"`},
		{"with quotes", NullString{String: `say "hi"`, Valid: true}, `"say \"hi\""`},
		{"with backslash", NullString{String: `a\b`, Valid: true}, `"a\\b"`},
		{"with newline", NullString{String: "a\nb", Valid: true}, `"a\nb"`},
		{"with tab", NullString{String: "a\tb", Valid: true}, `"a\tb"`},
		{"with carriage return", NullString{String: "a\rb", Valid: true}, `"a\rb"`},
		{"with ctrl char", NullString{String: "a\x01b", Valid: true}, "\"a\\u0001b\""},
		{"html angle brackets", NullString{String: "<script>", Valid: true}, "\"\\u003cscript\\u003e\""},
		{"html ampersand", NullString{String: "a&b", Valid: true}, "\"a\\u0026b\""},
		{"unicode", NullString{String: "café", Valid: true}, `"café"`},
		{"emoji", NullString{String: "hello 🌍", Valid: true}, `"hello 🌍"`},
		{"line separator", NullString{String: "a\u2028b", Valid: true}, "\"a\\u2028b\""},
		{"paragraph separator", NullString{String: "a\u2029b", Valid: true}, "\"a\\u2029b\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.val.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			// Verify against json.Marshal for valid strings
			if tt.val.Valid {
				expected, err := json.Marshal(tt.val.String)
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != string(expected) {
					t.Fatalf("mismatch with json.Marshal:\ngot:  %s\nwant: %s", got, expected)
				}
			}
			if string(got) != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNullTime_MarshalJSON(t *testing.T) {
	ts := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	tests := []struct {
		name string
		val  NullTime
		want string
	}{
		{"null", NullTime{}, "null"},
		{"valid", NullTime{Time: ts, Valid: true}, `"2024-03-15T10:30:00Z"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.val.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNullableMarshalJSON_RoundTrip(t *testing.T) {
	type record struct {
		Name  NullString `json:"name"`
		Age   NullInt64  `json:"age"`
		Born  NullTime   `json:"born"`
		Admin NullBool   `json:"admin"`
	}

	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	original := record{
		Name:  NullString{String: "Alice", Valid: true},
		Age:   NullInt64{Int64: 30, Valid: true},
		Born:  NullTime{Time: ts, Valid: true},
		Admin: NullBool{Bool: true, Valid: true},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded record
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Name.String != "Alice" || !decoded.Name.Valid {
		t.Fatalf("Name mismatch: %+v", decoded.Name)
	}
	if decoded.Age.Int64 != 30 || !decoded.Age.Valid {
		t.Fatalf("Age mismatch: %+v", decoded.Age)
	}
	if !decoded.Born.Time.Equal(ts) || !decoded.Born.Valid {
		t.Fatalf("Born mismatch: %+v", decoded.Born)
	}
	if decoded.Admin.Bool != true || !decoded.Admin.Valid {
		t.Fatalf("Admin mismatch: %+v", decoded.Admin)
	}

	// Round-trip with nulls
	nullRecord := record{}
	data, err = json.Marshal(nullRecord)
	if err != nil {
		t.Fatal(err)
	}
	var decodedNull record
	if err := json.Unmarshal(data, &decodedNull); err != nil {
		t.Fatal(err)
	}
	if decodedNull.Name.Valid || decodedNull.Age.Valid || decodedNull.Born.Valid || decodedNull.Admin.Valid {
		t.Fatal("expected all fields to be null")
	}
}

func BenchmarkNullString_MarshalJSON(b *testing.B) {
	ns := NullString{String: "hello world", Valid: true}
	b.ReportAllocs()
	for b.Loop() {
		ns.MarshalJSON()
	}
}

func BenchmarkNullString_MarshalJSON_Escaping(b *testing.B) {
	ns := NullString{String: `hello "world" <script>`, Valid: true}
	b.ReportAllocs()
	for b.Loop() {
		ns.MarshalJSON()
	}
}

func BenchmarkNullTime_MarshalJSON(b *testing.B) {
	ns := NullTime{Time: time.Now(), Valid: true}
	b.ReportAllocs()
	for b.Loop() {
		ns.MarshalJSON()
	}
}

func BenchmarkNullInt64_MarshalJSON(b *testing.B) {
	ns := NullInt64{Int64: 1234567890, Valid: true}
	b.ReportAllocs()
	for b.Loop() {
		ns.MarshalJSON()
	}
}

func BenchmarkNullBool_MarshalJSON(b *testing.B) {
	ns := NullBool{Bool: true, Valid: true}
	b.ReportAllocs()
	for b.Loop() {
		ns.MarshalJSON()
	}
}
