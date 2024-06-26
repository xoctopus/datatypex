package datatypex_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/datatypex"
)

func TestPassword(t *testing.T) {
	password := Password("any")
	NewWithT(t).Expect(password.String()).To(Equal("any"))
	NewWithT(t).Expect(password.SecurityString()).To(Equal(MaskedPassword))
}
