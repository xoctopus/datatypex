package types_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	. "github.com/sincospro/types"
)

func TestAddress_MarshalText(t *testing.T) {
	cases := []*struct {
		Name   string
		Addr   *Address
		Expect string
	}{
		{
			Name: "Asset",
			Addr: &Address{
				Group: "avatar",
				Key:   "filename",
				Ext:   "png",
			},
			Expect: "asset://avatar/filename.png",
		}, {
			Name: "HttpFileURI",
			Addr: &Address{
				URL: "https://demo.com/avatar/filename.png",
			},
			Expect: "https://demo.com/avatar/filename.png",
		}, {
			Name: "WithoutExtension",
			Addr: &Address{
				Group: "avatar",
				Key:   "filename",
			},
			Expect: "asset://avatar/filename",
		}, {
			Name: "LocalFile",
			Addr: &Address{
				URL: "file:///AbsPath/To/Your/Local/File.ext",
			},
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
	cases := []*struct {
		Name   string
		Input  string
		OutVal *Address
		OutErr error
	}{
		{
			Name:  "Asset",
			Input: "asset://avatar/filename.png",
			OutVal: &Address{
				URL:   "",
				Group: "avatar",
				Key:   "filename",
				Ext:   "png",
			},
		}, {
			Name:  "HttpFileURL",
			Input: "https://group.com/avatar/filename.png",
			OutVal: &Address{
				URL: "https://group.com/avatar/filename.png",
			},
		}, {
			Name:  "LocalFile",
			Input: "file:///AbsPath/To/Your/Local/File.ext",
			OutVal: &Address{
				URL: "file:///AbsPath/To/Your/Local/File.ext",
			},
		}, {
			Name:   "InvalidURI",
			Input:  "http://foo.com/ctl\x7f",
			OutErr: errors.New("should parse failed"),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			v := &Address{}
			err := v.Scan(c.Input)
			if c.OutErr != nil {
				NewWithT(t).Expect(err).NotTo(BeNil())
			} else {
				NewWithT(t).Expect(v).To(Equal(c.OutVal))
			}

		})
	}
	NewWithT(t).Expect((&Address{}).DataType("")).To(Equal("varchar(1024)"))
}
