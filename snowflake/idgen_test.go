package snowflake_test

import (
	"testing"
	"time"

	"github.com/xoctopus/datatypex"
	. "github.com/xoctopus/datatypex/snowflake"
	. "github.com/xoctopus/datatypex/snowflake/internal"
)

var (
	base, _ = time.ParseInLocation(datatypex.DefaultTimestampLayout, "2025-05-21T00:00:00.000UTC", time.UTC)
	g1      = NewIDGen(3)
	g2      = NewSnowflake(2, 1, base, 10, 12)
)

func Benchmark(b *testing.B) {
	b.Run("Generator", func(b *testing.B) {
		for range b.N {
			_ = g1.ID()
		}
	})

	b.Run("Snowflake", func(b *testing.B) {
		for range b.N {
			_ = g2.ID()
		}
	})
}
