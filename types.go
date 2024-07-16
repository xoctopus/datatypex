package datatypex

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"reflect"
)

// SqlValue can convert between sql value and describe sql datatype
type SqlValue interface {
	driver.Value
	sql.Scanner
	DataType(engine string) string
}

type TextValue interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

type JSONValue interface {
	json.Unmarshaler
	json.Marshaler
}

// Assertions
var (
	_ SqlValue = (*Address)(nil)
	_ SqlValue = (*Timestamp)(nil)
	_ SqlValue = (*UUID)(nil)
)

// Reflects
var (
	RtTextUnmarshaller = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	RtTextMarshaller   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// Interfaces
type (
	DefaultSetter    interface{ SetDefault() }
	CanBeZero        interface{ Zero() bool }
	Stringer         interface{ String() string }
	SecurityStringer interface{ SecurityString() string }
)
