package datatypex_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex"
)

func TestBinary(t *testing.T) {
	bytes, err := Binary("917").MarshalText()
	Expect(t, err, Succeed())
	Expect(t, bytes, Equal([]byte("917")))
	Expect(t, (&Binary{}).UnmarshalText(bytes), Succeed())
}
