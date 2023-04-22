package workflow

import "fmt"

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
