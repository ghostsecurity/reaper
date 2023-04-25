package transmission

import "encoding/json"

var _ Transmission = (*String)(nil)

type String string

func NewString(s string) *String {
	str := String(s)
	return &str
}

func (t *String) Type() Type {
	return NewType(TypeString, InternalTypeNone)
}

func (t *String) String() string {
	if t == nil {
		return ""
	}
	return string(*t)
}

func (t *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *String) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = String(v)
	return nil
}
