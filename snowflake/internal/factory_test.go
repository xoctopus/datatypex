package internal_test

import (
	"flag"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/datatypex"
	. "github.com/xoctopus/datatypex/snowflake/internal"
)

var (
	base_, _ = time.Parse(datatypex.DefaultTimestampLayout, "2025-05-21T00:00:00.000Z")
	// benchN as unit for generate factories for benchmarking
	benchN int
)

func init() {
	flag.IntVar(&benchN, "unit", 1, "set unit to run benchmark, default to run `Benchmark` without skip gap")
}

func TestFactory(t *testing.T) {
	t.Run("NewFactory", func(t *testing.T) {
		t.Run("AfterBase", func(t *testing.T) {
			ExpectPanic[error](t, func() {
				NewFactory(1, time.Now().Add(time.Minute), 4, 4)
			}, ErrorContains("the base timestamp MUST before now."))

		})
		t.Run("MissingBits", func(t *testing.T) {
			ExpectPanic[error](t, func() {
				NewFactory(1, base_, 4)
			}, ErrorContains("worker bits and sequence bits MUST be assigned"))

		})
		t.Run("InvalidBits", func(t *testing.T) {
			ExpectPanic[error](t, func() {
				NewFactory(1, base_, 32, 1)
			}, ErrorContains("worker bits and sequence bits MUST be less than 32 and timestamp bits MUST be greater than 0."))
		})
		t.Run("InvalidUnit", func(t *testing.T) {
			ExpectPanic[error](t, func() {
				NewFactory(0, base_, 4, 4)
			}, ErrorContains("unit MUST be greater than 0"))
		})
		t.Run("TooShortTimestamp", func(t *testing.T) {
			ExpectPanic[error](t, func() {
				NewFactory(5, base_, 20, 20)
			}, ErrorContains("factory MUST be able to generate continuously for 10 years or longer from now"))
		})
	})

	t.Run("Gaps", func(t *testing.T) {
		for _, unit := range []int{1, 5, 10, 30} {
			f := NewFactory(unit, base_, 4, 4)
			for i := range int64(5) {
				gaps := f.Gaps(base_.Add(time.Duration(i*int64(unit)) * time.Millisecond))
				Expect(t, gaps, Equal(i+f.Gap0()))
			}
		}
	})

	t.Run("Mask", func(t *testing.T) {
		f := NewFactory(1, base_, 4, 4)
		Expect(t, f.Mask(0xF), Equal(uint32(0xF)))
		Expect(t, f.Mask(0x7), Equal(uint32(0x7)))
		Expect(t, f.Mask(0x72), Equal(uint32(0x2)))
	})

	t.Run("Next", func(t *testing.T) {
		for _, unit := range []int{1, 5, 10, 30} {
			f := NewFactory(unit, base_, 4, 4)
			for _, n := range []int64{1, 2, 5, 6, 8} {
				start := time.Now()
				_, _ = f.Elapsed(), f.Next(n)
				sub := int64(time.Since(start)) / int64(unit) / int64(time.Millisecond)
				// sub := int64(time.Now().Sub(start)) / int64(unit) / int64(time.Millisecond)
				Expect(t, n-1 <= sub && sub <= n+1, BeTrue())
			}
		}
	})
}
