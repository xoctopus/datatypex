package snowflake

import (
	"time"

	"github.com/xoctopus/datatypex"
	"github.com/xoctopus/datatypex/snowflake/internal"
)

type IDGen interface {
	ID() int64
}

var base, _ = time.ParseInLocation(datatypex.DefaultTimestampLayout, "2025-05-21T00:00:00.000UTC", time.UTC)

func NewIDGen(worker uint32) IDGen {
	return internal.NewSnowflake(worker, 1, base, 10, 12)
}
