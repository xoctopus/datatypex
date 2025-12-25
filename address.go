package datatypex

import (
	"database/sql/driver"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

func ParseAddress(text string) (*Address, error) {
	u, err := url.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address from %q, [cause:%w]", text, err)
	}
	a := &Address{}

	if u.Scheme != AddressSchemeName {
		a.url = u.String()
		return a, nil
	}
	a.group = u.Hostname()
	if len(u.Path) > 0 {
		a.key = u.Path[1:]
		if idx := strings.LastIndex(u.Path, "."); idx != -1 {
			a.key = u.Path[1:idx]
			a.ext = u.Path[idx+1:]
		}
	}
	return a, nil
}

func NewAddress(group, filename string) *Address {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	key := strings.TrimSuffix(filename, "."+ext)
	v := &Address{
		group: group,
		key:   key,
		ext:   ext,
	}
	v.url = fmt.Sprintf("%s://%s/%s", AddressSchemeName, v.group, v.key)
	if v.ext != "" {
		v.url += "." + v.ext
	}
	return v
}

type Address struct {
	url   string
	group string
	key   string
	ext   string
}

const AddressSchemeName = "asset"

func (a Address) String() string {
	if a.url != "" {
		return a.url
	}
	if a.group == "" && a.key == "" {
		return ""
	}
	u := fmt.Sprintf("%s://%s/%s", AddressSchemeName, a.group, a.key)
	if a.ext != "" {
		u += "." + a.ext
	}
	return u
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Address) UnmarshalText(text []byte) error {
	addr, err := ParseAddress(string(text))
	if err != nil {
		return err
	}
	*a = *addr
	return nil
}

func (a Address) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a *Address) Scan(src any) error {
	return a.UnmarshalText([]byte(src.(string)))
}

func (a Address) DBType(driver string) string {
	return "varchar(1024)"
}
