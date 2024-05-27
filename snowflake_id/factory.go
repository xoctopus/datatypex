package snowflake_id

import (
	"fmt"
	"time"

	"github.com/xoctopus/datatypex"
)

func NewSnowflakeFactory(bitsWorkerID, bitsSequence uint32, gap uint32, base time.Time) *SnowflakeFactory {
	if now := time.Now(); now.Before(base) {
		base = now
	}

	bitsTimestamp := 63 - (bitsWorkerID + bitsSequence)
	f := &SnowflakeFactory{
		bitsWorkerID: bitsWorkerID,
		maxWorkerID:  1<<bitsWorkerID - 1,
		bitsSequence: bitsSequence,
		maxSequence:  1<<bitsSequence - 1,
		maxTimestamp: 1<<bitsTimestamp - 1,
		unit:         time.Duration(gap) * time.Millisecond,
		base:         base,
	}
	f.baseUnits = f.Units(base)

	return f
}

// SnowflakeFactory to build snowflake id
// |---worker id---|---sequence---|---timestamp---|
type SnowflakeFactory struct {
	bitsWorkerID uint32
	maxWorkerID  uint32 // maxWorkerID  = 1<<bitsWorkerID - 1
	bitsSequence uint32
	maxSequence  uint32 // maxSequence  = 1<<bitsSequence - 1
	maxTimestamp uint64
	unit         time.Duration
	base         time.Time
	baseUnits    uint64
}

func (f *SnowflakeFactory) MaskSequence(seq uint32) uint32 { return seq & f.maxSequence }

// Units returns snowflake units since `t`
func (f *SnowflakeFactory) Units(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(f.unit))
}

// Elapsed units from now to base
func (f *SnowflakeFactory) Elapsed(ts time.Time) uint64 {
	units := f.Units(ts)
	if f.baseUnits > units {
		panic(f.errInvalidTimestamp(ts))
	}
	return units - f.baseUnits
}

func (f *SnowflakeFactory) Duration(ts time.Time, d time.Duration) time.Duration {
	return d*f.unit - time.Duration(ts.UnixNano())%f.unit
}

func (f *SnowflakeFactory) New(workerID uint32) (*Snowflake, error) {
	if workerID > f.maxWorkerID {
		workerID = f.maxWorkerID
	}
	return &Snowflake{
		f:        f,
		workerID: workerID,
	}, nil
}

func (f *SnowflakeFactory) Build(workerID, seq uint32, elapsed uint64) (uint64, error) {
	if elapsed > f.maxTimestamp {
		return 0, f.errOverMaxTimestamp(elapsed)
	}
	// |sign|    elapsed     |    sequence   |    worker id   |
	// |  1 | bits timestamp | bits sequence | bits worker id | (64bits)
	return elapsed<<(f.bitsSequence+f.bitsWorkerID) | uint64(seq)<<f.bitsWorkerID | uint64(workerID), nil
}

func (f *SnowflakeFactory) String() string {
	du := time.Duration(f.baseUnits+f.maxTimestamp) * f.unit
	endAt := time.Unix(
		int64(du/time.Second),
		int64(du%time.Second),
	)
	return fmt.Sprintf(
		"EndAt[%s]_MaxWorker[%d]_MaxSeq[%d]_MaxTs[%d]",
		endAt.Format(datatypex.DefaultTimestampLayout),
		f.maxWorkerID, f.maxSequence, f.maxTimestamp,
	)
}

func (f *SnowflakeFactory) errInvalidTimestamp(ts time.Time) error {
	return errBeforeBaseTime{f.base, ts}
}

func (f *SnowflakeFactory) errOverMaxTimestamp(elapsed uint64) error {
	return errOverMaxTimestamp{f.maxTimestamp, elapsed}
}
