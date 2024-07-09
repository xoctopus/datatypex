package datatypex

import (
	"fmt"
	"go/ast"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/textx"
)

type Endpoint struct {
	Scheme   string
	Host     string
	Port     uint16
	Base     string
	Path     string
	Username string
	Password Password
	Param    url.Values
}

func (e Endpoint) String() string {
	u := url.URL{
		Scheme:   e.Scheme,
		Host:     e.Hostname(),
		RawPath:  "",
		RawQuery: e.Param.Encode(),
	}

	if e.Base != "" {
		u.Path = "/" + e.Base
	}

	if e.Param != nil {
		u.RawQuery = e.Param.Encode()
	}

	if e.Username != "" || e.Password != "" {
		u.User = url.UserPassword(e.Username, e.Password.String())
	}

	s, err := url.QueryUnescape(u.String())
	must.NoErrorWrap(err, "failed to query unescape: %s", u.String())
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

func ParseEndpoint(text string) (*Endpoint, error) {
	u, err := url.Parse(text)
	if err != nil {
		return nil, err
	}

	ep := &Endpoint{}

	ep.Scheme = u.Scheme

	if q := u.Query(); len(q) > 0 {
		ep.Param = q
	}

	ep.Path = u.Path

	if len(u.Path) > 0 {
		ep.Base = strings.TrimPrefix(u.Path, "/")
	}

	ep.Host = u.Hostname()
	if port, err := strconv.ParseUint(u.Port(), 10, 16); err == nil {
		ep.Port = uint16(port)
	}

	if u.User != nil {
		ep.Username = u.User.Username()
		password, _ := u.User.Password()
		ep.Password = Password(password)
	}

	return ep, nil
}

func UnmarshalExtra(ext url.Values, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return ErrUnmarshalExtraNonPointer(rv.Type().String())
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrUnmarshalExtraNonStruct(rv.Type().String())
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i)
		fn := ft.Name

		if !ast.IsExported(fn) {
			continue
		}
		if tag, ok := ft.Tag.Lookup("name"); ok {
			key, _ := reflectx.ParseTagKeyAndFlags(tag)
			if key == "-" {
				continue
			}
			if key != "" {
				fn = key
			}
		}
		fv := rv.Field(i)
		value := ext.Get(fn)
		if value == "" {
			value = ft.Tag.Get("default")
		}
		if err := textx.UnmarshalText([]byte(value), fv.Addr()); err != nil {
			return err
		}
	}
	return nil
}
