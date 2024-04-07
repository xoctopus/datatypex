package snowflake_id_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/sincospro/types/snowflake_id"
)

func NewSnowflakeTestSuite(t *testing.T, n int) *SnowflakeTestSuite {
	return &SnowflakeTestSuite{
		T: t,
		N: n,
	}
}

type SnowflakeTestSuite struct {
	*testing.T
	N    int
	mm   sync.Map
	size atomic.Int64
}

func (s *SnowflakeTestSuite) ExpectN(n int) {
	NewWithT(s.T).Expect(s.size.Load()).To(Equal(int64(n)))
}

func (s *SnowflakeTestSuite) Run(sf *Snowflake) {
	for i := 1; i <= s.N; i++ {
		id, err := sf.ID()
		if err != nil {
			s.T.Log(err)
		}
		NewWithT(s.T).Expect(err).To(BeNil())
		s.mm.Store(id, struct{}{})
		s.size.Add(1)
	}
}

func BenchmarkSnowflake_ID(b *testing.B) {
	for i, vs := range cases {
		f := NewSnowflakeFactory(vs[0], vs[1], vs[2], base)
		s, err := f.New(1)
		if err != nil {
			b.Fatal(err)
		}
		name := fmt.Sprintf("%02d_%s", i, f.String())
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err = s.ID()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestSnowflake_ID(t *testing.T) {
	sf, err := NewSnowflakeFactory(10, 12, 1, base).New(294)
	NewWithT(t).Expect(err).To(BeNil())

	for i := 0; i < 10000; i++ {
		_, err := sf.ID()
		NewWithT(t).Expect(err).To(BeNil())
		time.Sleep(100 * time.Microsecond)
	}

	t.Run("Concurrent", func(t *testing.T) {
		suite := NewSnowflakeTestSuite(t, 100)
		g, err := NewSnowflake(1)
		NewWithT(t).Expect(err).To(BeNil())

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
}
