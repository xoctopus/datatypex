package internal_test

import (
	"flag"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/datatypex"
	. "github.com/xoctopus/datatypex/snowflake/internal"
)

var (
	base_, _ = time.ParseInLocation(datatypex.DefaultTimestampLayout, "2025-05-21T00:00:00.000UTC", time.UTC)
	// benchN as unit for generate factories for benchmarking
	benchN int
)

func init() {
	flag.IntVar(&benchN, "unit", 1, "set unit to run benchmark, default to run `Benchmark` without skip gap")
}

func TestFactory(t *testing.T) {
	t.Run("NewFactory", func(t *testing.T) {
		t.Run("AfterBase", func(t *testing.T) {
			defer func() {
				e := recover().(error).Error()
				sub := "the base timestamp MUST before now."
				NewWithT(t).Expect(e).To(ContainSubstring(sub))
			}()
			NewFactory(1, time.Now().Add(time.Minute), 4, 4)
		})
		t.Run("MissingBits", func(t *testing.T) {
			defer func() {
				e := recover().(error).Error()
				sub := "worker bits and sequence bits MUST be assigned"
				NewWithT(t).Expect(e).To(ContainSubstring(sub))
			}()
			NewFactory(1, base_, 4)
		})
		t.Run("InvalidBits", func(t *testing.T) {
			defer func() {
				e := recover().(error).Error()
				sub := "worker bits and sequence bits MUST be less than 32 and timestamp bits MUST be greater than 0."
				NewWithT(t).Expect(e).To(ContainSubstring(sub))
			}()
			NewFactory(1, base_, 32, 1)
		})
		t.Run("InvalidUnit", func(t *testing.T) {
			defer func() {
				e := recover().(error).Error()
				NewWithT(t).Expect(e).To(ContainSubstring("unit MUST be greater than 0"))
			}()
			NewFactory(0, base_, 4, 4)
		})
		t.Run("TooShortTimestamp", func(t *testing.T) {
			defer func() {
				e := recover().(error).Error()
				NewWithT(t).Expect(e).To(ContainSubstring("factory MUST be able to generate continuously for 10 years or longer from now"))
			}()
			NewFactory(5, base_, 20, 20)
		})
	})

	t.Run("Gaps", func(t *testing.T) {
		for _, unit := range []int{1, 5, 10, 30} {
			f := NewFactory(unit, base_, 4, 4)
			for i := range int64(5) {
				gaps := f.Gaps(base_.Add(time.Duration(i*int64(unit)) * time.Millisecond))
				NewWithT(t).Expect(gaps).To(Equal(i + f.Gap0()))
			}
		}
	})

	t.Run("Mask", func(t *testing.T) {
		f := NewFactory(1, base_, 4, 4)
		NewWithT(t).Expect(f.Mask(0xF)).To(Equal(uint32(0xF)))
		NewWithT(t).Expect(f.Mask(0x7)).To(Equal(uint32(0x7)))
		NewWithT(t).Expect(f.Mask(0x72)).To(Equal(uint32(0x2)))
	})

	t.Run("Next", func(t *testing.T) {
		for _, unit := range []int{1, 5, 10, 30} {
			f := NewFactory(unit, base_, 4, 4)
			for _, n := range []int64{1, 2, 5, 6, 8} {
				start := time.Now()
				_, _ = f.Elapsed(), f.Next(n)
				sub := int64(time.Now().Sub(start)) / int64(unit) / int64(time.Millisecond)
				NewWithT(t).Expect(n-1 <= sub && sub <= n+1).To(BeTrue())
			}
		}
	})
}
