package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const DbTimeLayout = "02.01.2006 15:04:05"

type FormattedTime struct{ time.Time }

func NewFormattedTime(t time.Time) FormattedTime { return FormattedTime{t.UTC()} }
func NowFormattedTime() FormattedTime            { return FormattedTime{time.Now().UTC()} }

func (ft FormattedTime) Value() (driver.Value, error) { return ft.Time.UTC().Format(DbTimeLayout), nil }

func (ft *FormattedTime) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		ft.Time = v.UTC()
		return nil
	case []byte:
		return ft.parseString(string(v))
	case string:
		return ft.parseString(v)
	default:
		return fmt.Errorf("unsupported Scan type %T", src)
	}
}

func (ft *FormattedTime) parseString(s string) error {
	layouts := []string{
		DbTimeLayout,
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999-07:00",
		"2006-01-02 15:04:05-07:00",
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, time.UTC); err == nil {
			ft.Time = t.UTC()
			return nil
		}
	}
	return fmt.Errorf("cannot parse time: %q", s)
}

func (ft FormattedTime) Equal(other FormattedTime) bool { return ft.Time.Equal(other.Time) }
func (ft FormattedTime) IsZero() bool                   { return ft.Time.IsZero() }

func (ft *FormattedTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ft.Time.UTC().Format(DbTimeLayout) + "\""), nil
}

func (ft *FormattedTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		ft.Time = time.Time{}
		return nil
	}
	t, err := time.ParseInLocation(DbTimeLayout, s, time.UTC)
	if err != nil {
		return err
	}
	ft.Time = t
	return nil
}
