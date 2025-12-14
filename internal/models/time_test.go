package models

import (
	"testing"
	"time"
)

func TestFormattedTime_ScanVariants(t *testing.T) {
	var ft FormattedTime
	// DB layout
	if err := ft.Scan("14.12.2025 20:01:13"); err != nil {
		t.Fatalf("db layout parse failed: %v", err)
	}
	// RFC3339Nano
	if err := ft.Scan("2025-12-14T19:58:11.495037Z"); err != nil {
		t.Fatalf("rfc3339 parse failed: %v", err)
	}
	// Space RFC3339 with fractional and tz
	if err := ft.Scan("2025-12-14 19:58:11.495037+00:00"); err != nil {
		t.Fatalf("space rfc3339 parse failed: %v", err)
	}
	// Space RFC3339 without fractional
	if err := ft.Scan("2025-12-14 19:58:11+00:00"); err != nil {
		t.Fatalf("space rfc3339 no frac parse failed: %v", err)
	}
}

func TestFormattedTime_JSON(t *testing.T) {
	ft := NewFormattedTime(time.Date(2025, 12, 14, 20, 1, 13, 0, time.UTC))
	b, err := ft.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "\"14.12.2025 20:01:13\"" {
		t.Fatalf("unexpected json: %s", string(b))
	}
}
