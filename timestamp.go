package datatypex

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// standard time zones
var (
	UTC = time.UTC
	CST = time.FixedZone("CST", 8*60*60) // CST Asian/China
	JST = time.FixedZone("JST", 9*60*60) // JST Asian/Japan
	SGT = time.FixedZone("SGT", 8*60*60) // SGT Asian/Singapore
)

var (
	TimestampZero     = Timestamp{time.Time{}}
	TimestampUnixZero = Timestamp{time.Unix(0, 0)}
)

// DefaultTimestampLayout default timestamp layout with millisecond precision
// and time zone
const DefaultTimestampLayout = "2006-01-02T15:04:05.000Z07"

var gDefaultTimeZone = UTC

func SetDefaultTimeZone(tz *time.Location) {
	gDefaultTimeZone = tz
}

func Now() Timestamp {
	return Timestamp{time.Now()}
}

func ParseTimestamp(s string) (Timestamp, error) {
	return ParseTimestampWithLayout(s, DefaultTimestampLayout)
}

func ParseTimestampWithLayout(input, layout string) (Timestamp, error) {
	t, err := time.Parse(layout, input)
	if err != nil {
		return TimestampUnixZero, err
	}
	return Timestamp{t}, nil
}

// openapi:strfmt date-time
type Timestamp struct{ time.Time }

func (Timestamp) DBType(driver string) string {
	switch driver {
	case "postgres", "mysql":
		return "bigint"
	default:
		return "integer"
	}
}

func (t *Timestamp) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return NewErrTimestampScanBytes(v)
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
		return NewErrTimestampScanInvalidInput(v)
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

// String time string with default timestamp layout
func (t Timestamp) String() string {
	if t.IsZero() {
		return ""
	}
	return t.In(gDefaultTimeZone).Format(DefaultTimestampLayout)
}

func (t Timestamp) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Timestamp) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*t = TimestampZero
		return nil
	}
	_t, err := ParseTimestamp(string(data))
	if err != nil {
		return err
	}
	*t = _t
	return nil
}

func (t Timestamp) IsZero() bool {
	return t.Time.IsZero() || t == TimestampZero || t == TimestampUnixZero
}
