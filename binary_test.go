package datatypes_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/sincospro/datatypes"
)

func TestBinary(t *testing.T) {
	bytes, err := Binary("917").MarshalText()
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(bytes).To(Equal([]byte("917")))

	NewWithT(t).Expect((&Binary{}).UnmarshalText(bytes)).To(BeNil())

}
