package types_test

import (
	. "github.com/onsi/gomega"
	. "github.com/sincospro/types"
	"strconv"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	ts := Now()

	// Timestamp.DataType
	NewWithT(t).Expect(ts.DataType("")).To(Equal("bigint"))

	// ParseTimestampFromString
	ts2, err := ParseTimestampFromString(ts.String())
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(ts.Equal(ts2)).To(BeTrue())

	// ParseTimestampFromStringWithLayout
	ts2, err = ParseTimestampFromStringWithLayout(ts.String(), DefaultTimestampLayout)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(ts.Equal(ts2)).To(BeTrue())
	ts2, err = ParseTimestampFromStringWithLayout(ts.String(), "invalid timestamp layout")
	NewWithT(t).Expect(err).NotTo(BeNil())

	// Equals
	t1 := Timestamp{Time: ts.Time}
	t2 := Timestamp{Time: ts.Time}
	NewWithT(t).Expect(t1.Equal(t2)).To(BeTrue())
	NewWithT(t).Expect(t1.EqualSeconds(t2)).To(BeTrue())
	NewWithT(t).Expect(t1.EqualMillionSeconds(t2)).To(BeTrue())
	NewWithT(t).Expect(t1.EqualMicroSeconds(t2)).To(BeTrue())

	// Scan
	t3 := Timestamp{}
	err = t3.Scan("unsupported value type")
	NewWithT(t).Expect(err).NotTo(BeNil())

	err = t3.Scan([]byte("unsupported value type"))
	NewWithT(t).Expect(err).NotTo(BeNil())

	err = t3.Scan(nil)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(t3).To(Equal(TimestampUnixZero))

	err = t3.Scan(strconv.AppendInt([]byte{}, ts.UnixMilli(), 10))
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(t3)

	t4 := Timestamp{}
	err = t4.Scan(int64(-1))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(t4).To(Equal(TimestampUnixZero))

	err = t4.Scan(ts.UnixMilli())
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(t3.Equal(t4)).To(BeTrue())
	t.Log(t4)

	// Value
	v, err := t4.Value()
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(v.(int64)).To(Equal(t4.UnixMilli()))
	v, err = TimestampZero.Value()
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(v.(int64)).To(Equal(int64(0)))

	NewWithT(t).Expect(TimestampZero.String()).To(Equal(""))
	NewWithT(t).Expect(TimestampUnixZero.String()).To(Equal(""))
	NewWithT(t).Expect(ts.String()).NotTo(Equal(""))

	// String
	t.Log(t4)
	t.Log(TimestampZero)
	t.Log(TimestampUnixZero)

	// MarshalText
	text, err := t4.MarshalText()
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(string(text))

	// UnmarshalText
	t5 := Timestamp{}
	NewWithT(t).Expect(t5.UnmarshalText(text)).To(BeNil())
	NewWithT(t).Expect(t5.Equal(t4)).To(BeTrue())

	NewWithT(t).Expect(t5.UnmarshalText([]byte("invalid timestamp string"))).NotTo(BeNil())

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
