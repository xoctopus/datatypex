package datatypex_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/datatypex"
)

var (
	cause = errors.New("initialize with error")
)

type TestInitializer struct{}

func (i *TestInitializer) Init() {}

type TestInitializerWithError struct{}

func (i *TestInitializerWithError) Init() error { return cause }

type TestInitializerByContext struct{}

func (i *TestInitializerByContext) Init(_ context.Context) {}

type TestInitializerByContextWithError struct{}

func (i *TestInitializerByContextWithError) Init(_ context.Context) error { return cause }

func TestInit(t *testing.T) {
	cases := []struct {
		value  any
		result error
	}{
		{&TestInitializer{}, nil},
		{&TestInitializerWithError{}, cause},
		{&TestInitializerByContext{}, nil},
		{&TestInitializerByContextWithError{}, cause},
		{struct{}{}, nil},
	}

	for _, v := range cases {
		r1 := datatypex.Init(v.value)
		r2 := datatypex.InitByContext(context.Background(), v.value)
		if v.result == nil {
			Expect(t, r1, BeNil[error]())
			Expect(t, r2, BeNil[error]())
		} else {
			Expect(t, r1, Equal(v.result))
			Expect(t, r2, Equal(v.result))
		}
	}
}
