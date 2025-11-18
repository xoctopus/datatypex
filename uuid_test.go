package datatypex_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/datatypex"
)

func TestUUID(t *testing.T) {
	uuid, err := datatypex.NewUUID()
	Expect(t, err, Succeed())

	t.Log(uuid.String())

	Expect(t, uuid.DBType("postgres"), Equal("uuid"))
	Expect(t, uuid.DBType("mysql"), Equal("varchar(36)"))
	Expect(t, uuid.DBType("other_engine"), Equal("text"))
}
