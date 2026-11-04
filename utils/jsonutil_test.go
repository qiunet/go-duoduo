package utils

import (
	"testing"
	"time"
)

type sample struct {
	Name string    `json:"name"`
	At   time.Time `json:"at"`
	Age  int       `json:"age"`
}

type sampleOmit struct {
	Name string    `json:"name"`
	At   time.Time `json:"at,omitempty"`
}

func TestMarshal_Nil(t *testing.T) {
	b, err := Marshal(nil)
	if err != nil {
		t.Fatalf("Marshal(nil) error: %v", err)
	}
	if b != nil {
		t.Fatalf("Marshal(nil) = %v, want nil bytes", b)
	}
}

func TestUnmarshal_EmptyAndNull(t *testing.T) {
	var s sample
	s.Name = "keep"
	if err := Unmarshal(nil, &s); err != nil {
		t.Fatalf("Unmarshal(nil): %v", err)
	}
	if s.Name != "keep" {
		t.Fatalf("Unmarshal(nil) should not change value, got %+v", s)
	}

	if err := Unmarshal([]byte("null"), &s); err != nil {
		t.Fatalf("Unmarshal(null): %v", err)
	}
	if s.Name != "keep" {
		t.Fatalf("Unmarshal(null) should not change value, got %+v", s)
	}
}

func TestToJSONString_Nil(t *testing.T) {
	got, err := ToJSONString(nil)
	if err != nil {
		t.Fatalf("ToJSONString(nil) error: %v", err)
	}
	if got != "" {
		t.Fatalf("ToJSONString(nil) = %q, want empty", got)
	}
}

func TestFromJSONString_EmptyAndNull(t *testing.T) {
	var s sample
	s.Name = "keep"
	if err := FromJSONString("", &s); err != nil {
		t.Fatalf("FromJSONString(\"\"): %v", err)
	}
	if s.Name != "keep" {
		t.Fatalf("empty string should not change value, got %+v", s)
	}

	if err := FromJSONString("null", &s); err != nil {
		t.Fatalf("FromJSONString(\"null\"): %v", err)
	}
	if s.Name != "keep" {
		t.Fatalf("null should not change value, got %+v", s)
	}
}

func TestToJSONString_Basic(t *testing.T) {
	at := time.Unix(1700000000, 0)
	got, err := ToJSONString(sample{Name: "alice", At: at, Age: 18})
	if err != nil {
		t.Fatalf("ToJSONString: %v", err)
	}
	want := `{"name":"alice","at":1700000000,"age":18}`
	if got != want {
		t.Fatalf("ToJSONString = %s, want %s", got, want)
	}
}

func TestToJSONString_ZeroTimeAsNull(t *testing.T) {
	got, err := ToJSONString(sample{Name: "bob"})
	if err != nil {
		t.Fatalf("ToJSONString: %v", err)
	}
	want := `{"name":"bob","at":null,"age":0}`
	if got != want {
		t.Fatalf("ToJSONString = %s, want %s", got, want)
	}
}

func TestToJSONString_OmitEmptyZeroTime(t *testing.T) {
	got, err := ToJSONString(sampleOmit{Name: "carol"})
	if err != nil {
		t.Fatalf("ToJSONString: %v", err)
	}
	want := `{"name":"carol"}`
	if got != want {
		t.Fatalf("ToJSONString = %s, want %s", got, want)
	}
}

func TestFromJSONString_UnixTimestamp(t *testing.T) {
	var s sample
	err := FromJSONString(`{"name":"dave","at":1700000000,"age":20}`, &s)
	if err != nil {
		t.Fatalf("FromJSONString: %v", err)
	}
	if s.Name != "dave" || s.Age != 20 {
		t.Fatalf("unexpected fields: %+v", s)
	}
	if !s.At.Equal(time.Unix(1700000000, 0)) {
		t.Fatalf("At = %v, want unix 1700000000", s.At)
	}
}

func TestFromJSONString_RFC3339(t *testing.T) {
	var s sample
	err := FromJSONString(`{"name":"eve","at":"2023-11-14T22:13:20Z","age":21}`, &s)
	if err != nil {
		t.Fatalf("FromJSONString: %v", err)
	}
	want := time.Date(2023, 11, 14, 22, 13, 20, 0, time.UTC)
	if !s.At.Equal(want) {
		t.Fatalf("At = %v, want %v", s.At, want)
	}
}

func TestFromJSONString_NullTime(t *testing.T) {
	var s sample
	err := FromJSONString(`{"name":"frank","at":null,"age":22}`, &s)
	if err != nil {
		t.Fatalf("FromJSONString: %v", err)
	}
	if !s.At.IsZero() {
		t.Fatalf("At should be zero, got %v", s.At)
	}
}

func TestRoundTrip(t *testing.T) {
	in := sample{
		Name: "grace",
		At:   time.Unix(1710000000, 123).UTC(),
		Age:  30,
	}
	str, err := ToJSONString(in)
	if err != nil {
		t.Fatalf("ToJSONString: %v", err)
	}

	var out sample
	if err := FromJSONString(str, &out); err != nil {
		t.Fatalf("FromJSONString: %v", err)
	}
	if out.Name != in.Name || out.Age != in.Age {
		t.Fatalf("round trip fields mismatch: %+v vs %+v", out, in)
	}
	// 只保留秒级精度
	if out.At.Unix() != in.At.Unix() {
		t.Fatalf("round trip time mismatch: %v vs %v", out.At, in.At)
	}
}

func TestMarshalUnmarshal_Bytes(t *testing.T) {
	in := sample{Name: "heidi", At: time.Unix(1600000000, 0), Age: 1}
	b, err := Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var out sample
	if err := Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if out.Name != in.Name || out.Age != in.Age || out.At.Unix() != in.At.Unix() {
		t.Fatalf("bytes round trip mismatch: %+v vs %+v", out, in)
	}
}
