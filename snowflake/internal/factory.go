package internal

import (
	"fmt"
	"math"
	mbits "math/bits"
	"time"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/datatypex"
)

type part struct {
	bits int
	max  int64
}

// NewFactory creates a snowflake id factory
//
// requirements:
// 1. MUST generate more than 10 sfids per second
// 2. MUST be able to generate continuously for 10 years or longer
//
// The `unit` helps mitigate clock rollback to some extent. When the system
// clock moves backward within the `unit` range, the computed factory gap
// remains consistent, avoiding conflicts.
func NewFactory(unit int, base time.Time, bits ...int) *Factory {
	must.BeTrueWrap(
		base.Before(time.Now()),
		"the base timestamp MUST before now.",
	)
	must.BeTrueWrap(
		len(bits) == 2,
		"worker bits and sequence bits MUST be assigned",
	)
	must.BeTrueWrap(
		bits[0] < 32 && bits[1] < 32 && 63-bits[0]-bits[1] > 0,
		"worker bits and sequence bits MUST be less than 32 and timestamp bits MUST be greater than 0.",
	)
	must.BeTrueWrap(
		unit > 0,
		"unit MUST be greater than 0",
	)
	bits = append(bits, 63-bits[0]-bits[1])
	f := &Factory{
		wid:  part{bits[0], 1<<bits[0] - 1},
		seq:  part{bits[1], 1<<bits[1] - 1},
		gap:  part{bits[2], 1<<bits[2] - 1},
		unit: time.Duration(unit) * time.Millisecond,
		base: base,
	}
	end := time.Now().Add(10 * 365 * 24 * 60 * 60 * time.Second)
	hi, lo := mbits.Mul64(uint64(f.gap.max), uint64(f.unit))
	ts := time.Unix(int64(lo)/int64(time.Second), int64(lo)%int64(time.Second))
	must.BeTrueWrap(
		hi > 0 || ts.Sub(end) > 0,
		"factory MUST be able to generate continuously for 10 years or longer from now",
	)
	if hi > 0 || lo > uint64(math.MaxInt64) {
		f.end = f.base.Add(time.Duration(math.MaxInt64))
	} else {
		f.end = f.base.Add(time.Duration(lo))
	}

	f.gap0 = f.Gaps(f.base)
	f.tag = fmt.Sprintf("0x%02X_0x%02X_0x%02X", unit, bits[1], bits[0])

	return f
}

type Factory struct {
	// tag factory build tag
	tag string
	// wid worker id builder info
	wid part
	// seq builder info
	seq part
	// gap is builder info
	gap part
	// unit is the smallest gap for generate a snowflake id
	unit time.Duration
	// base is the base timestamp for this factory
	base time.Time
	// gap0 is the gaps from 1970-00-00T00:00:00.000 to base
	gap0 int64
	// end is the end timestamp for this factory
	end time.Time
}

func (f *Factory) Tag() string { return f.tag }

func (f *Factory) Gap0() int64 { return f.gap0 }

func (f *Factory) Gaps(t time.Time) int64 { return t.UnixNano() / int64(f.unit) }

// Elapsed returns units from base to t
func (f *Factory) Elapsed() int64 {
	gaps := f.Gaps(time.Now())
	must.BeTrue(gaps >= f.gap0)
	return gaps - f.gap0
}

func (f *Factory) Mask(seq uint32) uint32 { return seq & uint32(f.seq.max) }

func (f *Factory) New(wid uint32) *Snowflake {
	must.BeTrue(wid <= uint32(f.wid.max))
	return &Snowflake{f: f, id: wid}
}

func (f *Factory) Build(worker, seq uint32, gap int64) int64 {
	must.BeTrueWrap(
		gap <= f.gap.max,
		"assigned elapsed %d is greater than max units of factory %d",
		gap, f.gap.max,
	)
	// |sign|    elapsed     |    sequence   |    worker id   |
	// |  1 | bits timestamp | bits sequence | bits worker id | (64bits)
	// return gap<<(f.seq.bits+f.wid.bits) | int64(seq)<<f.wid.bits | int64(worker)
	return int64(worker)<<(f.gap.bits+f.seq.bits) | int64(seq)<<f.gap.bits | gap
}

func (f *Factory) Next(n int64) int64 {
	if n > 0 {
		d := time.Duration(n)*f.unit - time.Duration(time.Now().UnixNano())%f.unit
		time.Sleep(d)
	}
	return f.Elapsed()
}

func (f *Factory) String() string {
	return fmt.Sprintf(
		"BASE[%s]_END[%s]_%s",
		f.base.In(datatypex.CST).Format(datatypex.DefaultTimestampLayout),
		f.end.In(datatypex.CST).Format(datatypex.DefaultTimestampLayout),
		f.Tag(),
	)
}

func (f *Factory) Unit() int { return int(f.unit / time.Millisecond) }

func (f *Factory) SeqBits() int { return f.seq.bits }

func (f *Factory) WorkerBits() int { return f.wid.bits }
