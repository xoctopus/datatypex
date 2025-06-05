package datatypex_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/datatypex"
)

func TestUUID(t *testing.T) {
	uuid, err := datatypex.NewUUID()
	NewWithT(t).Expect(err).To(BeNil())

	t.Log(uuid.String())

	NewWithT(t).Expect(uuid.DBType("postgres")).To(Equal("uuid"))
	NewWithT(t).Expect(uuid.DBType("mysql")).To(Equal("varchar(36)"))
	NewWithT(t).Expect(uuid.DBType("other_engine")).To(Equal("text"))
}
