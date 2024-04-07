package snowflake_id

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	errInvalidSystemClock = errors.New("invalid system clock")
)

func NewSnowflake(worker uint32) (*Snowflake, error) {
	start, _ := time.Parse(time.RFC3339, "2018-10-24T07:30:06Z")
	return NewSnowflakeFactory(10, 12, 1, start).New(worker)
}

type Snowflake struct {
	f        *SnowflakeFactory
	workerID uint32
	elapsed  uint64
	sequence uint32
	mtx      sync.Mutex
}

func (s *Snowflake) WorkerID() uint32 { return s.workerID }

func (s *Snowflake) ID() (uint64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	now := time.Now()

	elapsed := s.f.Elapsed(now)
	if s.elapsed < elapsed {
		s.elapsed = elapsed
		s.sequence = RandU32N(s.f.maxSequence)
		return s.f.Build(s.workerID, s.sequence, s.elapsed)
	}

	if s.elapsed > elapsed {
		return 0, errInvalidSystemClock
	}

	s.sequence = s.f.MaskSequence(s.sequence + 1)
	if s.sequence == 0 {
		s.elapsed = s.elapsed + 1
		time.Sleep(s.f.Duration(now, time.Duration(s.elapsed-elapsed)))
	}

	return s.f.Build(s.workerID, s.sequence, s.elapsed)
}
