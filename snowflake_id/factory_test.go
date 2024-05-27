package snowflake_id_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/datatypex/snowflake_id"
)

var (
	base, _ = time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	cases   = [][3]uint32{
		{11, 11, 1},
		{10, 12, 1},
		{8, 16, 1},
		{8, 8, 4},
		{12, 10, 1},
		{12, 12, 1},
		{16, 10, 1},
		{16, 8, 1},
		{16, 8, 5},
		{16, 8, 10},
		{32, 24, 10},
		{32, 24, 100},
		{32, 28, 8},
	}
)

func TestSnowflakeFactory_Build(t *testing.T) {
	f := NewSnowflakeFactory(32, 28, 10, base)

	t.Log(f.String())

	id, err := f.Build(1, 2, 6)
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(id)

	id, err = f.Build(1, 2, 7)
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(id)

	_, err = f.Build(1, 2, 8)
	NewWithT(t).Expect(err).NotTo(BeNil())
	t.Log(err)
}

func TestSnowflakeFactory_Elapsed(t *testing.T) {
	f := NewSnowflakeFactory(12, 12, 10, base)

	t.Log(f.String())

	t.Run("EqualBase", func(t *testing.T) {
		t.Log(f.Elapsed(base))
	})

	t.Run("Elapsed", func(t *testing.T) {
		ts, err := time.Parse(time.RFC3339, "2020-01-01T03:00:00Z")
		NewWithT(t).Expect(err).To(BeNil())
		t.Log(f.Elapsed(ts))
	})

	t.Run("CatchPanic", func(t *testing.T) {
		defer func() {
			_err := recover()
			NewWithT(t).Expect(_err).NotTo(BeNil())
			t.Log(_err)
		}()
		ts, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
		NewWithT(t).Expect(err).To(BeNil())
		f.Elapsed(ts)
	})
}

func TestSnowflakeFactory_String(t *testing.T) {
	for i, c := range cases {
		f := NewSnowflakeFactory(c[0], c[1], c[2], base)
		t.Logf("%02d %s", i, f)
	}
}

func TestSnowflakeFactory_New(t *testing.T) {
	f := NewSnowflakeFactory(10, 12, 1, base)
	sf, err := f.New(1)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(sf.WorkerID()).To(Equal(uint32(1)))
	sf, err = f.New(1 << 10)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(sf.WorkerID()).To(Equal(uint32(1<<10 - 1)))
}

func TestNewSnowflakeFactory(t *testing.T) {
	NewSnowflakeFactory(10, 12, 1, base)
	NewSnowflakeFactory(10, 12, 1, time.Now().Add(time.Minute))
}
