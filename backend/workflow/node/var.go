package node

import (
	"encoding/json"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
)

type VarStorage struct {
	inputs  Connectors
	outputs Connectors
	static  map[string]transmission.Transmission
}

func NewVarStorage(inputs, outputs Connectors, defaults map[string]transmission.Transmission) *VarStorage {
	return &VarStorage{
		inputs:  inputs,
		outputs: outputs,
		static:  defaults,
	}
}

type VarStorageM struct {
	Inputs  Connectors               `json:"inputs"`
	Outputs Connectors               `json:"outputs"`
	Static  map[string]TransmissionM `json:"static"`
}

type TransmissionM struct {
	ParentType uint32      `json:"type"`
	ChildType  uint32      `json:"internal"`
	Data       interface{} `json:"data"`
}

func (m *TransmissionM) Unpack() (transmission.Transmission, error) {
	data, err := json.Marshal(m.Data)
	if err != nil {
		return nil, err
	}
	t := transmission.NewType(transmission.ParentType(m.ParentType), transmission.InternalType(m.ChildType))
	return transmission.UnmarshalJSON(t, data)
}

func (s *VarStorage) AddStaticInputValue(key string, value transmission.Transmission) error {
	input, ok := s.GetInputs().FindByName(key)
	if !ok {
		return fmt.Errorf("unexpected input '%s'", key)
	}
	if err := transmission.NewType(input.Type, 0).Validate(value); err != nil {
		return fmt.Errorf("invalid value for input '%s': %s", key, err)
	}
	if s.static == nil {
		s.static = make(map[string]transmission.Transmission)
	}
	s.static[key] = value
	return nil
}

func (s *VarStorage) Pack() (*VarStorageM, error) {
	m := make(map[string]TransmissionM, len(s.static))
	for name, t := range s.static {
		data, err := t.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %v", name, err)
		}
		var any interface{}
		if err := json.Unmarshal(data, &any); err != nil {
			return nil, fmt.Errorf("unmarshal %s: %v", name, err)
		}
		m[name] = TransmissionM{
			ParentType: uint32(t.Type().Parent()),
			ChildType:  uint32(t.Type().Internal()),
			Data:       any,
		}
	}
	return &VarStorageM{
		Inputs:  s.inputs,
		Outputs: s.outputs,
		Static:  m,
	}, nil
}

func (s *VarStorage) MarshalJSON() ([]byte, error) {
	packed, err := s.Pack()
	if err != nil {
		return nil, err
	}
	return json.Marshal(packed)
}

func (m *VarStorageM) Unpack() (*VarStorage, error) {
	s := &VarStorage{
		inputs:  m.Inputs,
		outputs: m.Outputs,
		static:  make(map[string]transmission.Transmission, len(m.Static)),
	}
	for name, t := range m.Static {
		tt, err := t.Unpack()
		if err != nil {
			return nil, err
		}
		s.static[name] = tt
	}
	return s, nil
}

func (s *VarStorage) UnmarshalJSON(data []byte) error {
	var m VarStorageM
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	ss, err := m.Unpack()
	if err != nil {
		return err
	}
	*s = *ss
	return nil
}

func (s *VarStorage) GetInputs() Connectors {
	return s.inputs
}

func (s *VarStorage) GetOutputs() Connectors {
	return s.outputs
}

func (s *VarStorage) SetStaticInputValues(values map[string]transmission.Transmission) error {
	if err := s.validate(values); err != nil {
		return err
	}
	s.static = values
	return nil
}

func (s *VarStorage) validate(values map[string]transmission.Transmission) error {
	for key, value := range values {
		input, ok := s.GetInputs().FindByName(key)
		if !ok {
			return fmt.Errorf("unexpected input '%s'", key)
		}
		if err := transmission.NewType(input.Type, 0).Validate(value); err != nil {
			return fmt.Errorf("invalid value for input '%s': %s", key, err)
		}
	}
	return nil
}

func (s *VarStorage) Validate(params map[string]transmission.Transmission) error {
	if err := s.validate(s.static); err != nil {
		return fmt.Errorf("invalid static inputs: %s", err)
	}
	if err := s.validate(params); err != nil {
		return fmt.Errorf("invalid dynamic inputs: %s", err)
	}
	return nil
}

func (s *VarStorage) ReadValue(name string, dynamicInputs map[string]transmission.Transmission) (transmission.Transmission, error) {
	input, ok := s.inputs.FindByName(name)
	if !ok {
		return nil, fmt.Errorf("input '%s' not found", name)
	}
	if input.Linkable && dynamicInputs != nil {
		if val, ok := dynamicInputs[name]; ok {
			return val, nil
		}
	}
	if s.static != nil {
		if val, ok := s.static[name]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("input '%s' not found", name)
}

func (s *VarStorage) ReadInputString(name string, dynamicInputs map[string]transmission.Transmission) (string, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return "", err
	}
	if v, ok := val.(transmission.Stringer); ok {
		return v.String(), nil
	}
	return "", fmt.Errorf("input '%s' is not a string", name)
}

func (s *VarStorage) ReadInputInt(name string, dynamicInputs map[string]transmission.Transmission) (int, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return 0, err
	}
	if v, ok := val.(transmission.Inter); ok {
		return v.Int(), nil
	}
	return 0, fmt.Errorf("input '%s' is not an int", name)
}

func (s *VarStorage) ReadInputBool(name string, dynamicInputs map[string]transmission.Transmission) (bool, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return false, err
	}
	if v, ok := val.(transmission.Booler); ok {
		return v.Bool(), nil
	}
	return false, fmt.Errorf("input '%s' is not a boolean", name)
}

func (s *VarStorage) ReadInputList(name string, dynamicInputs map[string]transmission.Transmission) (transmission.Lister, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return nil, err
	}
	if v, ok := val.(transmission.Lister); ok {
		return v, nil
	}
	return nil, fmt.Errorf("input '%s' is not a list", name)
}

func (s *VarStorage) ReadInputRequest(name string, dynamicInputs map[string]transmission.Transmission) (*packaging.HttpRequest, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return nil, err
	}
	if v, ok := val.(transmission.Requester); ok {
		r := v.Request()
		return &r, nil
	}
	return nil, fmt.Errorf("input '%s' is not a request", name)
}

func (s *VarStorage) ReadInputResponse(name string, dynamicInputs map[string]transmission.Transmission) (*packaging.HttpResponse, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return nil, err
	}
	if v, ok := val.(transmission.Responser); ok {
		r := v.Response()
		return &r, nil
	}
	return nil, fmt.Errorf("input '%s' is not a response", name)
}

func (s *VarStorage) ReadInputMap(name string, dynamicInputs map[string]transmission.Transmission) (map[string]string, error) {
	val, err := s.ReadValue(name, dynamicInputs)
	if err != nil {
		return nil, err
	}
	if v, ok := val.(transmission.Mapper); ok {
		return v.Map(), nil
	}
	return nil, fmt.Errorf("input '%s' is not a map", name)
}

func (s *VarStorage) FindInput(name string) (Connector, bool) {
	return s.inputs.FindByName(name)
}

func (s *VarStorage) FindOutput(name string) (Connector, bool) {
	return s.outputs.FindByName(name)
}
