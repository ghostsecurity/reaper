package transmission

type Start struct {
}

func NewStart() *Start {
	return &Start{}
}

func (e *Start) Type() Type {
	return NewType(TypeStart, InternalTypeNone)
}

func (e *Start) MarshalJSON() ([]byte, error) {
	return []byte(`null`), nil
}

func (e *Start) UnmarshalJSON([]byte) error {
	return nil
}
