package datatypex

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"reflect"
)

type (
	DefaultSetter    interface{ SetDefault() }
	ZeroChecker      interface{ IsZero() bool }
	Stringer         interface{ String() string }
	SecurityStringer interface{ SecurityString() string }
)

var (
	TypeDefaultSetter    = reflect.TypeFor[DefaultSetter]()
	TypeZeroChecker      = reflect.TypeFor[ZeroChecker]()
	TypeStringer         = reflect.TypeFor[Stringer]()
	TypeSecurityStringer = reflect.TypeFor[SecurityStringer]()
)

// DBEngineType identifies rdb engine type, usually it is `postgres`, `mysql`,
// `sqlite` or `sqlite3`, etc.
type DBEngineType string

// DBValue can convert between rdb value and go value with description of rdb datatype
type DBValue interface {
	driver.Valuer
	sql.Scanner
	DBType(engine DBEngineType) string
}

type TextArshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type JSONArshaler interface {
	json.Marshaler
	json.Unmarshaler
}
