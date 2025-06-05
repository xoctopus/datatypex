package datatypex_test

import (
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	. "github.com/xoctopus/datatypex"
)

var (
	AsErrParseEndpointByURL *ErrParseEndpointByURL
)

func TestParseEndpoint(t *testing.T) {
	cases := map[string]struct {
		uri    string
		expect *Endpoint
	}{
		"STMPs": {
			uri: "stmps://mail.xxx.com:465",
			expect: &Endpoint{
				Scheme: "stmps",
				Host:   "mail.xxx.com",
				Port:   465,
			},
		},
		"Postgres": {
			uri: "postgres://username:password@hostname:5432/database_name?sslmode=disable",
			expect: &Endpoint{
				Scheme:   "postgres",
				Host:     "hostname",
				Username: "username",
				Password: "password",
				Port:     5432,
				Base:     "database_name",
				Param:    url.Values{"sslmode": {"disable"}},
			},
		},
		"HTTPs": {
			uri: "https://hostname/path/to/resource?page=1&q=go 语言",
			expect: &Endpoint{
				Scheme: "https",
				Host:   "hostname",
				Base:   "path/to/resource",
				Param:  url.Values{"q": {"go 语言"}, "page": {"1"}},
			},
		},
		"NoScheme": {
			uri: "//hostname:1234/path/to/resource",
			expect: &Endpoint{
				Host: "hostname",
				Port: 1234,
				Base: "path/to/resource",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ep, err := ParseEndpoint(c.uri)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(ep.String()).To(Equal(c.uri))
			NewWithT(t).Expect(ep.String()).To(Equal(c.expect.String()))
			NewWithT(t).Expect(ep.SecurityString()).To(Equal(c.expect.SecurityString()))
			NewWithT(t).Expect(ep.IsZero()).To(Equal(c.expect.IsZero()))
			NewWithT(t).Expect(ep.Hostname()).To(Equal(c.expect.Hostname()))

			text, err := c.expect.MarshalText()
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(text).To(Equal([]byte(c.uri)))

			err = ep.UnmarshalText(text)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(ep.String()).To(Equal(c.uri))
		})
	}

	t.Run("FailedToParseURL", func(t *testing.T) {
		input := "http://hostname:http/path/to/resource"
		_, err := ParseEndpoint(input)
		NewWithT(t).Expect(errors.As(err, &AsErrParseEndpointByURL)).To(BeTrue())
		err = (&Endpoint{}).UnmarshalText([]byte(input))
		NewWithT(t).Expect(errors.As(err, &AsErrParseEndpointByURL)).To(BeTrue())
	})
}
