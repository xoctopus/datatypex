package datatypex_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex"
)

func TestPassword(t *testing.T) {
	password := Password("any")

	Expect(t, password.String(), Equal("any"))
	Expect(t, password.SecurityString(), Equal(MaskedPassword))
}
