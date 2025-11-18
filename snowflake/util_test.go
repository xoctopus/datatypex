package snowflake_test

import (
	randv1 "math/rand"
	randv2 "math/rand/v2"
	"net"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex/snowflake"
)

func TestWorkerIDFromIP(t *testing.T) {
	for _, c := range []*struct {
		ipv4Addr string
		workerId uint32
	}{
		{"255.255.255.255", 65535},
		{"127.0.0.1", 1},
		{"", 0},
	} {
		addr := net.ParseIP(c.ipv4Addr)
		wid := WorkerIDFromIP(addr)
		Expect(t, wid, Equal(c.workerId))
	}
	_, err := WorkerIDFromLocalIP()
	Expect(t, err, Succeed())
}

func BenchmarkRand(b *testing.B) {
	b.Run("v1", func(b *testing.B) {
		for range b.N {
			_ = randv1.Uint32()
		}
	})

	b.Run("v1_r", func(b *testing.B) {
		r := randv1.New(randv1.NewSource(time.Now().UnixNano()))
		for range b.N {
			_ = r.Uint32()
		}
	})

	b.Run("v2", func(b *testing.B) {
		for range b.N {
			_ = randv2.Uint32()
		}

	})

	b.Run("v2_r", func(b *testing.B) {
		ts := time.Now().UnixNano()
		r := randv2.New(randv2.NewPCG(uint64(ts<<32), uint64(ts>>32)))
		for range b.N {
			_ = r.Uint32()
		}
	})
}
