package datatypex

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseEndpoint(text string) (*Endpoint, error) {
	u, err := url.Parse(text)
	if err != nil {
		return nil, NewErrParseEndpointByURL(text, err)
	}

	ep := &Endpoint{
		Scheme: u.Scheme,
		Param:  u.Query(),
		Base:   strings.TrimPrefix(u.Path, "/"),
	}

	ep.Host = u.Hostname()
	if u.Port() != "" {
		port, _ := strconv.ParseUint(u.Port(), 10, 16)
		ep.Port = uint16(port)
	}

	if u.User != nil {
		ep.Username = u.User.Username()
		password, _ := u.User.Password()
		ep.Password = Password(password)
	}

	return ep, nil
}

type Endpoint struct {
	Scheme   string
	Host     string
	Port     uint16
	Base     string
	Username string
	Password Password
	Param    url.Values
}

func (e Endpoint) String() string {
	u := &url.URL{
		Scheme:   e.Scheme,
		Host:     e.Hostname(),
		RawQuery: e.Param.Encode(),
	}
	if e.Base != "" {
		u.Path = "/" + e.Base
	}

	if e.Username != "" || e.Password != "" {
		u.User = url.UserPassword(e.Username, e.Password.String())
	}

	s, _ := url.QueryUnescape(u.String())
	return s
}

func (e Endpoint) SecurityString() string {
	if e.Password != "" {
		e.Password = Password(e.Password.SecurityString())
	}
	return e.String()
}

func (e Endpoint) IsZero() bool {
	return e.Host == ""
}

func (e Endpoint) Hostname() string {
	if e.Port == 0 {
		return e.Host
	}
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (e Endpoint) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Endpoint) UnmarshalText(text []byte) error {
	ep, err := ParseEndpoint(string(text))
	if err != nil {
		return err
	}
	*e = *ep
	return nil
}
