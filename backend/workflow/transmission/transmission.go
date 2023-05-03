package transmission

import (
	"encoding/json"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/packaging"
)

type Type uint64

func NewType(parent ParentType, internal InternalType) Type {
	return Type(uint64(parent)<<32 | uint64(internal))
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint64{
		uint64(t) >> 32,
		uint64(t) & 0xFFFFFFFF,
	})
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var v []uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf("invalid type: %v", v)
	}
	*t = Type(v[0]<<32 | v[1])
	return nil
}

func (t Type) Parent() ParentType {
	return ParentType(uint64(t) >> 32)
}

func (t Type) Internal() InternalType {
	return InternalType(uint64(t) & 0xFFFFFFFF)
}

type ParentType uint32

const (
	// public types
	TypeUnknown ParentType = 1 << iota
	TypeString
	TypeInt
	TypeMap
	TypeList
	TypeRequest
	TypeResponse
	TypeStart
	TypeBoolean
	TypeAny = 0
)

type Stringer interface{ String() string }
type Inter interface{ Int() int }
type Booler interface{ Bool() bool }
type Mapper interface{ Map() map[string]string }
type Lister interface {
	Next() (string, bool)
	Count() int
	Complete() bool
}
type Requester interface{ Request() packaging.HttpRequest }
type Responser interface{ Response() packaging.HttpResponse }

func (t Type) Validate(transmission Transmission) error {
	if t.Parent()&transmission.Type().Parent() != t.Parent() {
		return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
	}
	for _, concrete := range []ParentType{
		TypeString,
		TypeInt,
		TypeMap,
		TypeList,
		TypeRequest,
		TypeResponse,
	} {
		if t.Parent()&concrete == 0 {
			continue
		}
		switch concrete {
		case TypeString:
			if _, ok := transmission.(Stringer); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		case TypeInt:
			if _, ok := transmission.(Inter); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		case TypeMap:
			if _, ok := transmission.(Mapper); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		case TypeList:
			if _, ok := transmission.(Lister); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		case TypeRequest:
			if _, ok := transmission.(Requester); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		case TypeResponse:
			if _, ok := transmission.(Responser); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		case TypeBoolean:
			if _, ok := transmission.(Booler); !ok {
				return fmt.Errorf("invalid transmission type %q for type %q", transmission.Type(), t)
			}
		}
	}
	return nil
}

type InternalType uint32

const (
	InternalTypeNone InternalType = iota
	InternalTypeUnknown
	InternalTypeNumericRangeList
	InternalTypeWordlist
	InternalTypeCommaSeparatedList
)

func (t ParentType) Contains(other ParentType) bool {
	return t&other == other
}

type Transmission interface {
	Type() Type
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
