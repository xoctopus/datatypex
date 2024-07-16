package datatypex_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/datatypex"
)

func TestUUID(t *testing.T) {
	uuid, err := datatypex.NewUUID()

	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(uuid.DataType("postgres")).To(Equal("uuid"))
	NewWithT(t).Expect(uuid.DataType("mysql")).To(Equal("text"))
}
