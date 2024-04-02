package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"strconv"
	"time"
)

var (
	UTC               = time.UTC
	CST               = time.FixedZone("CST", 8*60*60)
	TimestampZero     = Timestamp{time.Time{}}
	TimestampUnixZero = Timestamp{time.Unix(0, 0)}
)

const (
	// DefaultTimestampLayout default timestamp layout
	DefaultTimestampLayout = "2006-01-02T15:04:05.000Z07:00"
)

// openapi:strfmt date-time
type Timestamp struct{ time.Time }

var (
	_ sql.Scanner              = (*Timestamp)(nil)
	_ driver.Valuer            = (*Timestamp)(nil)
	_ encoding.TextMarshaler   = (*Timestamp)(nil)
	_ encoding.TextUnmarshaler = (*Timestamp)(nil)
)

func (Timestamp) DataType(_ string) string {
	return "bigint"
}

func Now() Timestamp {
	return Timestamp{time.Now()}
}

func ParseTimestampFromString(s string) (Timestamp, error) {
	return ParseTimestampFromStringWithLayout(s, DefaultTimestampLayout)
}

func ParseTimestampFromStringWithLayout(input, layout string) (Timestamp, error) {
	t, err := time.Parse(layout, input)
	if err != nil {
		return TimestampUnixZero, err
	}
	return Timestamp{t}, nil
}

func (t *Timestamp) Equal(compared Timestamp) bool {
	return t.EqualMillionSeconds(compared)
}

func (t *Timestamp) EqualSeconds(compared Timestamp) bool {
	return t.Time.Unix() == compared.Time.Unix()
}

func (t *Timestamp) EqualMillionSeconds(compared Timestamp) bool {
	return t.Time.UnixMilli() == compared.Time.UnixMilli()
}

func (t *Timestamp) EqualMicroSeconds(compared Timestamp) bool {
	return t.Time.UnixMicro() == compared.Time.UnixMicro()
}

func (t *Timestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("sql.Scan() strfmt.Timestamp from: %#v failed: %s", v, err.Error())
		}
		*t = Timestamp{time.Unix(n/1000, n%1000*1e6)}
	case int64:
		if v < 0 {
			*t = TimestampUnixZero
		} else {
			*t = Timestamp{time.Unix(v/1000, v%1000*1e6)}
		}
	case nil:
		*t = TimestampUnixZero
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Timestamp from: %#v", v)
	}
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	ts := t.UnixMilli()
	if ts <= 0 {
		return int64(0), nil
	}
	return ts, nil
}

func (t Timestamp) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DefaultTimestampLayout)
}

func (t Timestamp) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Timestamp) UnmarshalText(data []byte) error {
	_t, err := ParseTimestampFromString(string(data))
	if err != nil {
		return err
	}
	*t = _t
	return nil
}

func (t Timestamp) IsZero() bool {
	if t.Time.IsZero() {
		return true
	}
	unix := t.Unix()
	return unix == 0 || unix == TimestampZero.Unix()
}
