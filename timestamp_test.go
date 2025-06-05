package datatypex_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/datatypex"
)

func ExampleTimestamp_String() {
	ts, _ := ParseTimestamp("2022-10-24T07:30:00.000CST")
	fmt.Println("SYS:", ts.String()) // default use time.Local

	SetDefaultTimeZone(UTC)
	ts, _ = ParseTimestamp("2022-10-24T07:30:00.000CST")
	fmt.Println("UTC:", ts.String())

	SetDefaultTimeZone(CST)
	ts, _ = ParseTimestamp("2022-10-24T07:30:00.000CST")
	fmt.Println("CST:", ts.String())

	SetDefaultTimeZone(JST)
	ts, _ = ParseTimestamp("2022-10-24T07:30:00.000CST")
	fmt.Println("JST:", ts.String())

	SetDefaultTimeZone(SGT)
	ts, _ = ParseTimestamp("2022-10-24T07:30:00.000CST")
	fmt.Println("SGT:", ts.String())

	// Output:
	// SYS: 2022-10-24T07:30:00.000CST
	// UTC: 2022-10-23T23:30:00.000UTC
	// CST: 2022-10-24T07:30:00.000CST
	// JST: 2022-10-24T08:30:00.000JST
	// SGT: 2022-10-24T07:30:00.000SGT
}

var (
	AsErrTimestampScanBytes        *ErrTimestampScanBytes
	AsErrTimestampScanInvalidInput *ErrTimestampScanInvalidInput
)

func TestTimestamp(t *testing.T) {
	SetDefaultTimeZone(CST)
	t.Run("ParseTimestampWithLayout", func(t *testing.T) {
		ts := Now()
		parsed, err := ParseTimestampWithLayout(ts.String(), DefaultTimestampLayout)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(ts.String()).To(Equal(parsed.String()))
		t.Run("InvalidLayout", func(t *testing.T) {
			_, err = ParseTimestampWithLayout(ts.String(), time.RFC822)
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
	})
	t.Run("DialectType", func(t *testing.T) {
		ts := Now()
		NewWithT(t).Expect(ts.DBType("postgres")).To(Equal("bigint"))
		NewWithT(t).Expect(ts.DBType("sqlite")).To(Equal("integer"))
	})
	t.Run("Value", func(t *testing.T) {
		ts, err := ParseTimestamp("1970-01-01T00:00:01.234UTC")
		NewWithT(t).Expect(err).To(BeNil())
		v, err := ts.Value()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v).To(Equal(int64(1234)))

		t.Run("InvalidOrZeroTimestamp", func(t *testing.T) {
			v, err = TimestampZero.Value()
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(v).To(Equal(int64(0)))
		})
	})
	t.Run("Scan", func(t *testing.T) {
		for _, v := range []*struct {
			input  any
			err    any
			result Timestamp
		}{
			{[]byte("1970-01-01T08:00:01.234CST"), AsErrTimestampScanBytes, Timestamp{}},
			{[]byte("1234"), nil, Timestamp{Time: time.UnixMilli(1234)}},
			{int64(-1), nil, TimestampUnixZero},
			{int64(0), nil, TimestampUnixZero},
			{int64(1234), nil, Timestamp{Time: time.UnixMilli(1234)}},
			{nil, nil, TimestampUnixZero},
			{"abc", AsErrTimestampScanInvalidInput, Timestamp{}},
		} {
			ts := Now()
			err := ts.Scan(v.input)
			if v.err != nil {
				NewWithT(t).Expect(errors.As(err, &v.err)).To(BeTrue())
			} else {
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(ts == v.result).To(BeTrue())
			}
		}
	})
	t.Run("String", func(t *testing.T) {
		NewWithT(t).Expect(TimestampUnixZero.String()).To(Equal(""))
		ts := Timestamp{Time: time.UnixMilli(1234)}
		NewWithT(t).Expect(ts.String()).To(Equal("1970-01-01T08:00:01.234CST"))
	})
	t.Run("TextArshaler", func(t *testing.T) {
		ts := Now()
		data, err := ts.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(data).To(Equal([]byte(ts.String())))

		err = ts.UnmarshalText([]byte("1970-01-01T00:00:01.234UTC"))
		NewWithT(t).Expect(err).To(BeNil())

		err = ts.UnmarshalText([]byte("1970-01-01T00:00:01.234"))
		NewWithT(t).Expect(err).NotTo(BeNil())

		err = ts.UnmarshalText([]byte(""))
		NewWithT(t).Expect(err).To(BeNil())
	})
	t.Run("IsZero", func(t *testing.T) {
		NewWithT(t).Expect(TimestampZero.IsZero()).To(BeTrue())
		NewWithT(t).Expect(TimestampUnixZero.IsZero()).To(BeTrue())
		ts := Timestamp{Time: time.UnixMicro(0)}
		NewWithT(t).Expect(ts.IsZero()).To(BeTrue())
		NewWithT(t).Expect(Now().IsZero()).To(BeFalse())
	})
}
