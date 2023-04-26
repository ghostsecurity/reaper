package transmission

import (
	"encoding/json"
)

var _ Transmission = (*Int)(nil)

type Int int

func NewInt(i int) *Int {
	i2 := Int(i)
	return &i2
}

func (i *Int) Type() Type {
	return NewType(TypeInt, InternalTypeNone)
}

func (i *Int) Int() int {
	if i == nil {
		return 0
	}
	return int(*i)
}

func (i *Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int())
}

func (i *Int) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = Int(v)
	return nil
}
