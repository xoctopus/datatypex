package internal

import (
	"sync"
	"time"

	"github.com/xoctopus/x/misc/must"
)

func NewSnowflake(worker uint32, unit int, base time.Time, w, s int) *Snowflake {
	return NewFactory(unit, base, w, s).New(worker)
}

type Snowflake struct {
	f   *Factory
	id  uint32
	seq uint32
	gap int64
	mtx sync.Mutex
}

func (s *Snowflake) WorkerID() uint32 { return s.id }

func (s *Snowflake) Tag() string { return s.f.Tag() }

func (s *Snowflake) ID() int64 {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	gap := s.f.Elapsed()
	if s.gap < gap {
		s.gap = gap
		// s.seq = _rand.Uint32N(uint32(s.f.seq.max))
		s.seq = 1
		return s.f.Build(s.id, s.seq, s.gap)
	}

	if s.gap > gap {
		gap = s.f.Elapsed()
		must.BeTrueF(s.gap <= gap, "invalid system clock, clock moved backwards")
	}

	// in same gap. if reached max sequence, need to wait next gap to make sure
	// generated in ascending order
	s.seq = s.f.Mask(s.seq + 1)
	if s.seq == 0 {
		s.gap++
		s.f.Next(s.gap - gap)
	}
	return s.f.Build(s.id, s.seq, s.gap)
}
