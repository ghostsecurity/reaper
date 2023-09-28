package transmission

import (
	"encoding/json"
	"fmt"
)

func UnmarshalJSON(t Type, data json.RawMessage) (Transmission, error) {
	var target Transmission
	switch t.Parent() {
	case TypeString:
		target = new(String)
	case TypeInt:
		target = new(Int)
	case TypeMap:
		target = new(Map)
	case TypeList:
		switch t.Internal() {
		case InternalTypeNumericRangeList:
			target = new(numericRangeIterator)
		case InternalTypeWordlist:
			target = new(wordlistIterator)
		case InternalTypeCommaSeparatedList:
			target = new(csvIterator)
		default:
			return nil, fmt.Errorf("unknown internal list type %q", t.Internal())
		}
	case TypeRequest:
		target = new(Request)
	case TypeResponse:
		target = new(Response)
	case TypeRequest | TypeResponse:
		target = new(RequestResponsePair)
	case TypeBoolean:
		target = new(Boolean)
	case TypeStart:
		target = new(Start)
	case TypeChoice:
		target = new(Choice)
	default:
		return nil, fmt.Errorf("unknown type %q", t)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}
	return target, nil
}
