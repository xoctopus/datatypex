package datatypex

import (
	"fmt"
	"reflect"
)

func typ(v any) string {
	switch x := v.(type) {
	case reflect.Value:
		return x.Type().String()
	default:
		return fmt.Sprintf("%T", x)
	}
}

func NewErrParseAddressByURL(input string, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrParseAddressByURL{input, cause}
}

type ErrParseAddressByURL struct {
	input string
	cause error
}

func (e *ErrParseAddressByURL) Error() string {
	return fmt.Sprintf("failed to parse address from `%s`: [%+v]", e.input, e.cause)
}

func (e *ErrParseAddressByURL) Unwrap() error {
	return e.cause
}

func NewErrParseEndpointByURL(input string, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrParseEndpointByURL{input, cause}
}

type ErrParseEndpointByURL struct {
	input string
	cause error
}

func (e *ErrParseEndpointByURL) Error() string {
	return fmt.Sprintf("failed to parse endpoint from `%s`: [%+v]", e.input, e.cause)
}

func (e *ErrParseEndpointByURL) Unwrap() error {
	return e.cause
}

func NewErrTimestampScanBytes(input []byte) error {
	return &ErrTimestampScanBytes{input: input}
}

type ErrTimestampScanBytes struct {
	input []byte
}

func (e *ErrTimestampScanBytes) Error() string {
	return fmt.Sprintf("failed to sql.Scan() strfmt.Timestamp from: %#v", e.input)
}

func NewErrTimestampScanInvalidInput(got any) error {
	return &ErrTimestampScanInvalidInput{got: typ(got)}
}

type ErrTimestampScanInvalidInput struct {
	got string
}

func (e *ErrTimestampScanInvalidInput) Error() string {
	return "invalid sql.Scan() strfmt.Timestamp input type. expect []byte, int64 or a nil value, but got " + e.got
}

func NewErrUnmarshalSFID(input []byte, cause error) error {
	if cause == nil {
		return nil
	}
	return &ErrUnmarshalSFID{input, cause}
}

type ErrUnmarshalSFID struct {
	input []byte
	cause error
}

func (e *ErrUnmarshalSFID) Error() string {
	return fmt.Sprintf("failed to unmarshal SFID from `%s`: [%+v]", e.input, e.cause)
}

func (e *ErrUnmarshalSFID) Unwrap() error {
	return e.cause
}
