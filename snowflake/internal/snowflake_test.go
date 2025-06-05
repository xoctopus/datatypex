package internal_test

import (
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/xhd2015/xgo/runtime/mock"

	. "github.com/xoctopus/datatypex/snowflake/internal"
)

func BenchmarkSnowflake_ID(b *testing.B) {
	var (
		factories []*Factory
		worker    int
		seq       int
	)

	gen := func() {
		var f *Factory
		defer func() {
			if recover() == nil {
				factories = append(factories, f)
			}
		}()
		f = NewFactory(benchN, base_, worker, seq)
	}

	// The generated factories can support at least 64 workers and generating up
	// to 128 sfids per gap.
	for bits := 16; bits <= 62; bits++ {
		for worker = 6; worker <= bits; worker++ {
			seq = bits - worker
			if worker <= 6 || seq <= 8 {
				continue
			}
			gen()
		}
	}
	sort.Slice(factories, func(i, j int) bool {
		_i, _j := factories[i], factories[j]
		if _i.SeqBits() > _j.SeqBits() {
			return false
		}
		if _i.SeqBits() < _j.SeqBits() {
			return true
		}
		return _i.WorkerBits() < _j.WorkerBits()
	})

	for _, f := range factories {
		b.Run(f.Tag(), func(b *testing.B) {
			s := f.New(1)
			defer func() {
				if e := recover(); e != nil {
					b.Log(recover())
				}
			}()
			for range b.N {
				_ = s.ID()
			}
		})
	}
}

type SnowflakeTestSuite struct {
	*testing.T
	N    int
	m    sync.Map
	size atomic.Int64
}

func NewSnowflakeTestSuite(t *testing.T, n int) *SnowflakeTestSuite {
	return &SnowflakeTestSuite{T: t, N: n}
}

func (s *SnowflakeTestSuite) ExpectN(n int) {
	NewWithT(s.T).Expect(s.size.Load()).To(Equal(int64(n)))
}

func (s *SnowflakeTestSuite) Run(sf *Snowflake) {
	for range s.N {
		id := sf.ID()
		s.m.Store(id, struct{}{})
		s.size.Add(1)
	}
}

func TestSnowflake_ID(t *testing.T) {
	gap, worker, seq := 1, 10, 12
	f := NewFactory(1, base_, 10, 12)
	NewWithT(t).Expect(f.Unit()).To(Equal(gap))
	NewWithT(t).Expect(f.SeqBits()).To(Equal(seq))
	NewWithT(t).Expect(f.WorkerBits()).To(Equal(worker))

	g1 := f.New(1)
	g2 := NewSnowflake(1, 1, base_, 10, 12)
	NewWithT(t).Expect(g1.WorkerID()).To(Equal(g2.WorkerID()))
	NewWithT(t).Expect(g1.Tag()).To(Equal(g2.Tag()))

	t.Run(f.String()+"_1x", func(t *testing.T) {
		g := NewFactory(1, base_, 4, 4).New(1)

		for i := 0; i < 10000; i++ {
			func() {
				defer func() {
					NewWithT(t).Expect(recover()).To(BeNil())
				}()
				_ = g.ID()
			}()
		}
	})

	t.Run(f.Tag()+"_1000x", func(t *testing.T) {
		suite := NewSnowflakeTestSuite(t, 1000)
		g := f.New(1)

		con := 1000
		wg := &sync.WaitGroup{}

		for i := 0; i < con; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				suite.Run(g)
			}()
		}

		wg.Wait()
		suite.ExpectN(suite.N * con)
	})

	t.Run(f.Tag()+"InvalidClock", func(t *testing.T) {
		g1.ID()
		now := time.Now()

		defer func() {
			NewWithT(t).Expect(recover().(string)).To(Equal("invalid system clock, clock moved backwards"))
		}()

		mock.Patch(time.Now, func() time.Time { return now.Add(0 - 10*time.Second) })
		g1.ID()
	})
}
