package datatypex

import (
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

type TextArshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type JSONArshaler interface {
	json.Marshaler
	json.Unmarshaler
}
