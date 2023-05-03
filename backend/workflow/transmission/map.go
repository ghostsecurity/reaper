package transmission

import "encoding/json"

var _ Transmission = (*Map)(nil)

type Map map[string]string

func NewMap(data map[string]string) *Map {
	m := Map(data)
	return &m
}

func (t *Map) Type() Type {
	return NewType(TypeMap, InternalTypeNone)
}

func (t *Map) Map() map[string]string {
	if t == nil {
		return nil
	}
	return map[string]string(*t)
}

func (t *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Map())
}

func (t *Map) UnmarshalJSON(data []byte) error {
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = Map(v)
	return nil
}
