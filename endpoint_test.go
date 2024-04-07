package types_test

import (
	"net/url"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/sincospro/types"
)

func TestEndpoint(t *testing.T) {
	cases := map[string]struct {
		ep  *Endpoint
		url string
	}{
		"STMPs": {
			ep: &Endpoint{
				Scheme: "stmps",
				Host:   "stmps.mail.xxx.com",
				Port:   465,
			},
			url: "stmps://stmps.mail.xxx.com:465",
		},
		"Postgres": {
			ep: &Endpoint{
				Scheme:   "postgres",
				Host:     "hostname",
				Username: "username",
				Password: "password",
				Port:     5432,
				Path:     "/database_name",
				Base:     "database_name",
				Param:    url.Values{"sslmode": {"disable"}},
			},
			url: "postgres://username:password@hostname:5432/database_name?sslmode=disable",
		},
		"NoScheme": {
			ep: &Endpoint{
				Scheme:   "",
				Host:     "hostname",
				Username: "username",
				Password: "password",
				Port:     5432,
				Path:     "/database_name",
				Base:     "database_name",
				Param:    url.Values{"sslmode": {"disable"}},
			},
			url: "//username:password@hostname:5432/database_name?sslmode=disable",
		},
		"HostOnly": {
			ep: &Endpoint{
				Scheme: "https",
				Host:   "host",
				Base:   "path",
				Path:   "/path",
			},
			url: "https://host/path",
		},
		"FailedToParse1": {
			url: string([]byte{0x7f}),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.ep != nil {
				NewWithT(t).Expect(c.ep.String()).To(Equal(c.url))
				text, err := c.ep.MarshalText()
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(string(text)).To(Equal(c.url))
			}

			parsed, err1 := ParseEndpoint(c.url)
			unmarshaled := &Endpoint{}
			err2 := unmarshaled.UnmarshalText([]byte(c.url))

			if c.ep != nil {
				NewWithT(t).Expect(err1).To(BeNil())
				NewWithT(t).Expect(err2).To(BeNil())
				NewWithT(t).Expect(*parsed).To(Equal(*c.ep))
				NewWithT(t).Expect(*unmarshaled).To(Equal(*c.ep))
			} else {
				NewWithT(t).Expect(err1).NotTo(BeNil())
				NewWithT(t).Expect(err2).NotTo(BeNil())
			}
		})
	}

	t.Run("SecurityString", func(t *testing.T) {
		t.Log(cases["Postgres"].ep.SecurityString())
	})

	t.Run("IsZero", func(t *testing.T) {
		NewWithT(t).Expect((&Endpoint{}).IsZero()).To(BeTrue())
	})

	/*
		t.Run("UnmarshalExtra", func(t *testing.T) {
			opt := struct {
				ConnectTimeout Duration `name:"connectTimeout" default:"10s"`
				ReadTimeout    Duration `name:"readTimeout" default:"10s"`
				WriteTimeout   Duration `name:"writeTimeout" default:"10s"`
				IdleTimeout    Duration `name:"idleTimeout" default:"240s"`
				MaxActive      int      `name:"maxActive" default:"5"`
				MaxIdle        int      `name:"maxIdle" default:"3"`
				DB             int      `name:"db" default:"10"`
			}{}

			err := UnmarshalExtra(url.Values{}, &opt)
			NewWithT(t).Expect(err).To(BeNil())
			// spew.Dump(opt)
		})
	*/
}

func TestEndpoint_UnmarshalText(t *testing.T) {

}

func TestEndpoint_MarshalText(t *testing.T) {

}
