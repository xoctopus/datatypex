package datatypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"reflect"
)

// SqlValue can convert between sql value and describe sql datatype
type SqlValue interface {
	driver.Value
	sql.Scanner
	DataType(engine string) string
}

// Assertions
var (
	_ SqlValue = (*Address)(nil)
	_ SqlValue = (*Timestamp)(nil)
)

type ErrUnmarshalExtraNonPointer string

func (e ErrUnmarshalExtraNonPointer) Error() string {
	return "non-pointer value `" + string(e) + "` is not supported"
}

type ErrUnmarshalExtraNonStruct string

func (e ErrUnmarshalExtraNonStruct) Error() string {
	return "non-struct value `" + string(e) + "` is not supported"
}

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
