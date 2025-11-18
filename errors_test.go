package datatypex_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex"
)

type ErrorWrapper interface {
	Unwrap() error
}

func TestErrors(t *testing.T) {
	input := "input"
	cause := errors.New("cause")

	for _, e := range []error{
		NewErrParseAddressByURL(input, cause),
		NewErrParseEndpointByURL(input, cause),
		NewErrTimestampScanBytes([]byte(input)),
		NewErrTimestampScanInvalidInput(1),
		NewErrTimestampScanInvalidInput(reflect.ValueOf(1)),
		NewErrUnmarshalSFID([]byte(input), cause),
	} {
		Expect(t, e, NotBeNil[error]())
		t.Log(e)
		if u, ok := e.(ErrorWrapper); ok {
			Expect(t, u.Unwrap(), Equal(cause))
		}
	}

	for _, e := range []error{
		NewErrParseAddressByURL(input, nil),
		NewErrParseEndpointByURL(input, nil),
		NewErrUnmarshalSFID(nil, nil),
	} {
		Expect(t, e, BeNil[error]())
	}
}
