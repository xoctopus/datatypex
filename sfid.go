package datatypex

import "strconv"

// Number.MAX_SAFE_INTEGER JavaScript (2^53 – 1)

type SFID uint64

func (sf SFID) MarshalText() ([]byte, error) {
	return []byte(sf.String()), nil
}

func (sf *SFID) UnmarshalText(data []byte) error {
	str := string(data)
	if len(str) == 0 {
		*sf = 0
		return nil
	}
	u, err := strconv.ParseUint(str, 10, 64)
	*sf = SFID(u)
	return NewErrUnmarshalSFID(data, err)
}

func (sf SFID) String() string {
	return strconv.FormatUint(uint64(sf), 10)
}

func NewSFIDs(vs ...uint64) SFIDs {
	ids := make(SFIDs, len(vs))
	for i, v := range vs {
		ids[i] = SFID(v)
	}
	return ids
}

type SFIDs []SFID

func (sfs SFIDs) ToUint64() (integers []uint64) {
	for _, sf := range sfs {
		integers = append(integers, uint64(sf))
	}
	return
}
