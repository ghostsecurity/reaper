package transmission

import "encoding/json"

type Boolean struct {
	value bool
}

func NewBoolean(value bool) *Boolean {
	return &Boolean{
		value: value,
	}
}

func (e *Boolean) Type() Type {
	return NewType(TypeBoolean, InternalTypeNone)
}

func (e *Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value)
}

func (e *Boolean) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.value)
}

func (e *Boolean) Bool() bool {
	return e.value
}
