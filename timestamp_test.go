package datatypex_test

import (
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/datatypex"
)

func TestTimestamp(t *testing.T) {
	ts := Now()

	t.Run("DataType", func(t *testing.T) {
		NewWithT(t).Expect(ts.DataType("")).To(Equal("bigint"))
	})

	t.Run("ParseTimestampFromString", func(t *testing.T) {
		ts2, err := ParseTimestampFromString(ts.String())
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(ts.Equal(ts2)).To(BeTrue())
	})

	t.Run("ParseTimestampFromStringWithLayout", func(t *testing.T) {
		ts2, err := ParseTimestampFromStringWithLayout(ts.String(), DefaultTimestampLayout)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(ts.Equal(ts2)).To(BeTrue())
		ts2, err = ParseTimestampFromStringWithLayout(ts.String(), "invalid timestamp layout")
		NewWithT(t).Expect(err).NotTo(BeNil())
		t.Log(ts2)
	})

	t.Run("Equals", func(t *testing.T) {
		t1 := Timestamp{Time: ts.Time}
		t2 := Timestamp{Time: ts.Time}
		NewWithT(t).Expect(t1.Equal(t2)).To(BeTrue())
		NewWithT(t).Expect(t1.EqualSeconds(t2)).To(BeTrue())
		NewWithT(t).Expect(t1.EqualMillionSeconds(t2)).To(BeTrue())
		NewWithT(t).Expect(t1.EqualMicroSeconds(t2)).To(BeTrue())
	})

	t.Run("Scan", func(t *testing.T) {
		t3 := Timestamp{}
		t.Run("UnsupportedValueType", func(t *testing.T) {
			err := t3.Scan("unsupported value type")
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
		t.Run("ScanFromBytes", func(t *testing.T) {
			t.Run("InvalidBytesValue", func(t *testing.T) {
				err := t3.Scan([]byte("invalid value"))
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
			err := t3.Scan(strconv.AppendInt([]byte{}, ts.UnixMilli(), 10))
			NewWithT(t).Expect(err).To(BeNil())
			t.Log(t3)
		})
		t.Run("ScanFromInt", func(t *testing.T) {
			t.Run("Zero", func(t *testing.T) {
				err := t3.Scan(nil)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(t3).To(Equal(TimestampUnixZero))
			})
			t4 := Timestamp{}
			err := t4.Scan(int64(-1))
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(t3).To(Equal(TimestampUnixZero))

			err = t4.Scan(ts.UnixMilli())
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(ts.Equal(t4)).To(BeTrue())
			t.Log(t4)
		})
	})

	t.Run("Value", func(t *testing.T) {
		t4 := Timestamp{}

		NewWithT(t).Expect(t4.Scan(ts.UnixMilli())).To(BeNil())

		v, err := t4.Value()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v.(int64)).To(Equal(t4.UnixMilli()))

		v, err = TimestampZero.Value()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v.(int64)).To(Equal(int64(0)))
	})
	// Value

	t.Run("String", func(t *testing.T) {
		NewWithT(t).Expect(TimestampZero.String()).To(Equal(""))
		NewWithT(t).Expect(TimestampUnixZero.String()).To(Equal(""))
		NewWithT(t).Expect(ts.String()).NotTo(Equal(""))
		t.Log(ts)
		t.Log(TimestampZero)
		t.Log(TimestampUnixZero)
	})

	t.Run("MarshalText", func(t *testing.T) {
		text, err := ts.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		t.Log(string(text))
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		text, err := ts.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())

		t5 := Timestamp{}
		NewWithT(t).Expect(t5.UnmarshalText(text)).To(BeNil())
		NewWithT(t).Expect(t5.Equal(ts)).To(BeTrue())
		NewWithT(t).Expect(t5.UnmarshalText([]byte("invalid timestamp string"))).NotTo(BeNil())
	})

	// std.Time tests
	now := time.Now()
	shadow, err := time.Parse(time.RFC3339, now.Format(time.RFC3339))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(now.Equal(shadow)).To(BeFalse())

	shadow2, err := time.Parse(time.RFC3339Nano, now.Format(time.RFC3339Nano))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(now.Equal(shadow2)).To(BeTrue())

	t6 := time.Time{}
	t.Log(t6.IsZero())
	t7 := time.Unix(0, 0)
	t.Log(t7.IsZero())
	t.Log(TimestampZero.Unix())
	t.Log(TimestampUnixZero.Unix())
}
