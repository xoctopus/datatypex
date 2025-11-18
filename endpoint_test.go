package datatypex_test

import (
	"net/url"
	"testing"

	"github.com/pkg/errors"
	. "github.com/xoctopus/x/testx"

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
			Expect(t, err, Succeed())
			Expect(t, ep.String(), Equal(c.uri))
			Expect(t, ep.String(), Equal(c.expect.String()))
			Expect(t, ep.SecurityString(), Equal(c.expect.SecurityString()))
			Expect(t, ep.IsZero(), Equal(c.expect.IsZero()))
			Expect(t, ep.Hostname(), Equal(c.expect.Hostname()))

			text, err := c.expect.MarshalText()
			Expect(t, err, Succeed())
			Expect(t, text, Equal([]byte(c.uri)))

			err = ep.UnmarshalText(text)
			Expect(t, err, Succeed())
			Expect(t, ep.String(), Equal(c.uri))
		})
	}

	t.Run("FailedToParseURL", func(t *testing.T) {
		input := "http://hostname:http/path/to/resource"
		_, err := ParseEndpoint(input)
		Expect(t, errors.As(err, &AsErrParseEndpointByURL), BeTrue())
		err = (&Endpoint{}).UnmarshalText([]byte(input))
		Expect(t, errors.As(err, &AsErrParseEndpointByURL), BeTrue())
	})
}
