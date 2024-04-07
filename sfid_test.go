package datatypes_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/sincospro/datatypes"
)

func TestSFID(t *testing.T) {
	sfid := SFID(100)

	text, err := sfid.MarshalText()
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(string(text)).To(Equal("100"))

	err = sfid.UnmarshalText([]byte("101"))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(sfid).To(Equal(SFID(101)))

	err = sfid.UnmarshalText([]byte{})
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(sfid).To(Equal(SFID(0)))

	err = sfid.UnmarshalText([]byte("not number"))
	NewWithT(t).Expect(err).NotTo(BeNil())
	NewWithT(t).Expect(sfid).To(Equal(SFID(0)))
}

func TestSFIDs(t *testing.T) {
	sfids := SFIDs{1, 2, 3}
	NewWithT(t).Expect(sfids.ToUint64()).To(Equal([]uint64{1, 2, 3}))
}
