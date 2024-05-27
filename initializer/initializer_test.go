package initializer_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/xoctopus/datatypex/initializer"
)

var (
	errInitializerWithError          = errors.New("initialize with error")
	errInitializerByContextWithError = errors.New("initialize by context with error")
)

type Initializer struct{}

func (i *Initializer) Init() {}

type InitializerWithError struct{}

func (i *InitializerWithError) Init() error {
	return errInitializerWithError
}

type InitializerByContext struct{}

func (i *InitializerByContext) Init(_ context.Context) {}

type InitializerByContextWithError struct{}

func (i *InitializerByContextWithError) Init(_ context.Context) error {
	return errInitializerByContextWithError
}

func TestInit(t *testing.T) {
	results := []error{
		nil,
		errInitializerWithError,
		nil,
		errInitializerByContextWithError,
		nil,
	}
	for i, v := range []any{
		&Initializer{},
		&InitializerWithError{},
		&InitializerByContext{},
		&InitializerByContextWithError{},
		&struct{}{},
	} {
		if results[i] == nil {
			NewWithT(t).Expect(initializer.Init(v)).To(BeNil())
		} else {
			NewWithT(t).Expect(initializer.Init(v)).To(Equal(results[i]))
		}
	}
}
