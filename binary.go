package datatypex

type Binary []byte

func (d Binary) MarshalText() ([]byte, error) { return d, nil }

func (d *Binary) UnmarshalText(data []byte) (err error) {
	*d = Binary(data)
	return
}
