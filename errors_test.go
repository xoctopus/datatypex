package datatypex_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

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
		NewWithT(t).Expect(e).NotTo(BeNil())
		t.Log(e)
		if u, ok := e.(ErrorWrapper); ok {
			NewWithT(t).Expect(u.Unwrap()).To(Equal(cause))
		}
	}

	for _, e := range []error{
		NewErrParseAddressByURL(input, nil),
		NewErrParseEndpointByURL(input, nil),
		NewErrUnmarshalSFID(nil, nil),
	} {
		NewWithT(t).Expect(e).To(BeNil())
	}
}
