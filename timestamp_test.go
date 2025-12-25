package datatypex_test

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex"
)

func ExampleTimestamp_String() {
	input := "2022-10-24T07:30:00.000+08"

	SetDefaultTimeZone(UTC)
	ts, _ := ParseTimestamp(input)
	fmt.Println("UTC:", ts.String())

	SetDefaultTimeZone(CST)
	ts, _ = ParseTimestamp(input)
	fmt.Println("CST:", ts.String())

	SetDefaultTimeZone(JST)
	ts, _ = ParseTimestamp(input)
	fmt.Println("JST:", ts.String())

	SetDefaultTimeZone(SGT)
	ts, _ = ParseTimestamp(input)
	fmt.Println("SGT:", ts.String())

	// Output:
	// UTC: 2022-10-23T23:30:00.000Z
	// CST: 2022-10-24T07:30:00.000+08
	// JST: 2022-10-24T08:30:00.000+09
	// SGT: 2022-10-24T07:30:00.000+08
}

func TestTimestamp(t *testing.T) {
	SetDefaultTimeZone(CST)
	t.Run("ParseTimestampWithLayout", func(t *testing.T) {
		ts := Now()
		parsed, err := ParseTimestampWithLayout(ts.String(), DefaultTimestampLayout)
		Expect(t, err, Succeed())
		Expect(t, ts.String(), Equal(parsed.String()))
		t.Run("InvalidLayout", func(t *testing.T) {
			_, err = ParseTimestampWithLayout(ts.String(), time.RFC822)
			Expect(t, err, Failed())
		})
	})
	t.Run("DialectType", func(t *testing.T) {
		ts := Now()
		Expect(t, ts.DBType("postgres"), Equal("bigint"))
		Expect(t, ts.DBType("sqlite"), Equal("integer"))
	})
	t.Run("Value", func(t *testing.T) {
		ts, err := ParseTimestamp("1970-01-01T00:00:01.234Z")
		Expect(t, err, Succeed())
		v, err := ts.Value()
		Expect(t, err, Succeed())
		Expect(t, v, Equal[driver.Value](int64(1234)))

		t.Run("InvalidOrZeroTimestamp", func(t *testing.T) {
			v, err = TimestampZero.Value()
			Expect(t, err, Succeed())
			Expect(t, v, Equal[driver.Value](int64(0)))
		})
	})
	t.Run("Scan", func(t *testing.T) {
		for _, v := range []*struct {
			input  any
			failed bool
			result Timestamp
		}{
			{[]byte("1970-01-01T08:00:01.234CST"), true, Timestamp{}},
			{[]byte("1234"), false, Timestamp{Time: time.UnixMilli(1234)}},
			{int64(-1), false, TimestampUnixZero},
			{int64(0), false, TimestampUnixZero},
			{int64(1234), false, Timestamp{Time: time.UnixMilli(1234)}},
			{nil, false, TimestampUnixZero},
			{"abc", true, Timestamp{}},
		} {
			ts := Now()
			err := ts.Scan(v.input)
			if err != nil {
				Expect(t, v.failed, BeTrue())
			} else {
				Expect(t, err, Succeed())
				Expect(t, ts == v.result, BeTrue())
			}
		}
	})
	t.Run("String", func(t *testing.T) {
		Expect(t, TimestampUnixZero.String(), Equal(""))
		ts := Timestamp{Time: time.UnixMilli(1234)}
		Expect(t, ts.String(), Equal("1970-01-01T08:00:01.234+08"))
	})
	t.Run("TextArshaler", func(t *testing.T) {
		ts := Now()
		data, err := ts.MarshalText()
		Expect(t, err, Succeed())
		Expect(t, data, Equal([]byte(ts.String())))

		err = ts.UnmarshalText([]byte("1970-01-01T00:00:01.234+08"))
		Expect(t, err, Succeed())

		err = ts.UnmarshalText([]byte("1970-01-01T00:00:01.234CST"))
		Expect(t, err, Failed())

		err = ts.UnmarshalText([]byte(""))
		Expect(t, err, Succeed())
	})
	t.Run("IsZero", func(t *testing.T) {
		Expect(t, TimestampZero.IsZero(), BeTrue())
		Expect(t, TimestampUnixZero.IsZero(), BeTrue())
		ts := Timestamp{Time: time.UnixMicro(0)}
		Expect(t, ts.IsZero(), BeTrue())
		Expect(t, Now().IsZero(), BeFalse())
	})
}
