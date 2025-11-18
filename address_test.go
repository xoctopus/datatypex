package datatypex_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/xoctopus/x/misc/must"
	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/datatypex"
)

var (
	AsErrParseAddressByURL *ErrParseAddressByURL
)

func TestAddress_MarshalText(t *testing.T) {
	cases := []*struct {
		Name   string
		Addr   *Address
		Expect string
	}{
		{
			Name:   "Asset",
			Addr:   NewAddress("avatar", "filename.png"),
			Expect: "asset://avatar/filename.png",
		}, {
			Name:   "HttpFileURI",
			Addr:   must.NoErrorV(ParseAddress("https://demo.com/avatar/filename.png")),
			Expect: "https://demo.com/avatar/filename.png",
		}, {
			Name:   "WithoutExtension",
			Addr:   NewAddress("avatar", "filename"),
			Expect: "asset://avatar/filename",
		}, {
			Name:   "LocalFile",
			Addr:   must.NoErrorV(ParseAddress("file:///AbsPath/To/Your/Local/File.ext")),
			Expect: "file:///AbsPath/To/Your/Local/File.ext",
		}, {
			Name:   "Empty",
			Addr:   &Address{},
			Expect: "",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			bytes, err := c.Addr.Value()
			Expect(t, err, Succeed())
			Expect(t, bytes.(string), Equal(c.Expect))
			raw, err := c.Addr.MarshalText()
			Expect(t, err, Succeed())
			Expect(t, string(raw), Equal(c.Expect))
		})
	}

}

func TestAddress_UnmarshalText(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		OutVal *Address
		OutErr error
	}{
		{
			Name:   "Asset",
			Input:  "asset://avatar/filename.png",
			OutVal: NewAddress("avatar", "filename.png"),
		}, {
			Name:   "HttpFileURL",
			Input:  "https://group.com/avatar/filename.png",
			OutVal: must.NoErrorV(ParseAddress("https://group.com/avatar/filename.png")),
		}, {
			Name:   "LocalFile",
			Input:  "file:///AbsPath/To/Your/Local/File.ext",
			OutVal: must.NoErrorV(ParseAddress("file:///AbsPath/To/Your/Local/File.ext")),
		}, {
			Name:   "InvalidURI",
			Input:  "http://foo.com/ctl\x7f",
			OutErr: AsErrParseAddressByURL,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			v := &Address{}
			err := v.Scan(c.Input)
			if err != nil {
				Expect(t, errors.As(err, &c.OutErr), BeTrue())
			} else {
				Expect(t, v.String(), Equal(c.OutVal.String()))
			}
		})
	}
	Expect(t, (&Address{}).DBType(""), Equal("varchar(1024)"))
}
