package datatypes

import (
	"database/sql"
	"database/sql/driver"
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
