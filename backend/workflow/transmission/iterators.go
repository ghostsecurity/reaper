package transmission

import (
	"encoding/json"
	"fmt"
)

var _ Transmission = (*numericRangeIterator)(nil)

type numericRangeIterator struct {
	min     int
	max     int
	current int
}

func NewNumericRangeIterator(min, max int) *numericRangeIterator {
	if min > max {
		min, max = max, min
	}
	return &numericRangeIterator{
		min:     min,
		max:     max,
		current: min - 1,
	}
}

func (n *numericRangeIterator) Next() (string, bool) {
	if n.current >= n.max {
		return "", false
	}
	n.current++
	return fmt.Sprintf("%d", n.current), true
}

func (n *numericRangeIterator) Count() int {
	return (n.max - n.min) + 1
}

func (n *numericRangeIterator) Complete() bool {
	return n.current >= n.max
}

func (n *numericRangeIterator) Type() Type {
	return NewType(TypeList, InternalTypeNumericRangeList)
}

func (n *numericRangeIterator) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]int{
		n.min,
		n.max,
	})
}

func (n *numericRangeIterator) UnmarshalJSON(data []byte) error {
	var v [2]int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.min = v[0]
	n.max = v[1]
	n.current = v[0] - 1
	return nil
}
