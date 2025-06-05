package datatypex_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/xoctopus/x/misc/must"

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
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(bytes).To(Equal(c.Expect))
			raw, err := c.Addr.MarshalText()
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(string(raw)).To(Equal(c.Expect))
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
				NewWithT(t).Expect(errors.As(err, &c.OutErr)).NotTo(BeNil())
			} else {
				NewWithT(t).Expect(v.String()).To(Equal(c.OutVal.String()))
			}
		})
	}
	NewWithT(t).Expect((&Address{}).DBType("")).To(Equal("varchar(1024)"))
}
