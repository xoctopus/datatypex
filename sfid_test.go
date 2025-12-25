package datatypex_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex"
)

func TestSFID(t *testing.T) {
	sfid := SFID(100)

	text, err := sfid.MarshalText()
	Expect(t, err, BeNil[error]())
	Expect(t, string(text), Equal("100"))

	err = sfid.UnmarshalText([]byte("101"))
	Expect(t, err, BeNil[error]())
	Expect(t, sfid, Equal(SFID(101)))

	err = sfid.UnmarshalText([]byte{})
	Expect(t, err, BeNil[error]())
	Expect(t, sfid, Equal(SFID(0)))

	err = sfid.UnmarshalText([]byte("not number"))
	Expect(t, err, Failed())
	Expect(t, sfid, Equal(SFID(0)))
}

func TestSFIDs(t *testing.T) {
	sfids := NewSFIDs(1, 2, 3)
	Expect(t, sfids.ToUint64(), Equal([]uint64{1, 2, 3}))
}
