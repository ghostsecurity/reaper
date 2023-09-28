package transmission

import "encoding/json"

var _ Transmission = (*Choice)(nil)

type Choice struct {
	options map[string]string
	key     string
}

type jsonChoice struct {
	Options map[string]string `json:"options"`
	Key     string            `json:"key"`
}

func NewChoice(key string, options map[string]string) *Choice {
	return &Choice{
		options: options,
		key:     key,
	}
}

func (t *Choice) Type() Type {
	return NewType(TypeChoice, InternalTypeNone)
}

func (t *Choice) Choice() (string, map[string]string) {
	if t == nil {
		return "", nil
	}
	return t.key, t.options
}

func (t *Choice) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonChoice{
		Options: t.options,
		Key:     t.key,
	})
}

func (t *Choice) UnmarshalJSON(data []byte) error {
	var v jsonChoice
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.options = v.Options
	t.key = v.Key
	return nil
}
