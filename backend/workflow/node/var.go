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

func NewVarStorage(inputs, outputs Connectors) *VarStorage {
	return &VarStorage{
		inputs:  inputs,
		outputs: outputs,
		static:  make(map[string]transmission.Transmission),
	}
}

type mVarStorage struct {
	Inputs  Connectors               `json:"inputs"`
	Outputs Connectors               `json:"outputs"`
	Static  map[string]mTransmission `json:"static"`
}

type mTransmission struct {
	Type transmission.Type `json:"type"`
	Data json.RawMessage   `json:"data"`
}

func (s *VarStorage) MarshalJSON() ([]byte, error) {
	m := make(map[string]mTransmission, len(s.static))
	for name, t := range s.static {
		data, err := t.MarshalJSON()
		if err != nil {
			return nil, err
		}
		m[name] = mTransmission{
			Type: t.Type(),
			Data: data,
		}
	}
	return json.Marshal(mVarStorage{
		Inputs:  s.inputs,
		Outputs: s.outputs,
		Static:  m,
	})
}

func (s *VarStorage) UnmarshalJSON(data []byte) error {
	var m mVarStorage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	s.inputs = m.Inputs
	s.outputs = m.Outputs
	s.static = make(map[string]transmission.Transmission, len(m.Static))
	for name, t := range m.Static {
		tt, err := transmission.UnmarshalJSON(t.Type, t.Data)
		if err != nil {
			return err
		}
		s.static[name] = tt
	}
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
	if s.static != nil {
		if val, ok := s.static[name]; ok {
			return val, nil
		}
	}
	if dynamicInputs != nil {
		if val, ok := dynamicInputs[name]; ok {
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
