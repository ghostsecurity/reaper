package transmission

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

type wordlistIterator struct {
	filename string
	f        *os.File
	scanner  *bufio.Scanner
	complete bool
}

func NewWordlistIterator(filename string) *wordlistIterator {
	return &wordlistIterator{
		filename: filename,
	}
}

func (w *wordlistIterator) Next() (string, bool) {
	if w.complete {
		return "", false
	}
	if w.scanner == nil {
		f, err := os.Open(w.filename)
		if err != nil {
			return "", false
		}
		w.f = f
		w.scanner = bufio.NewScanner(f)
	}

	for {
		if !w.scanner.Scan() {
			break
		}
		text := w.scanner.Text()
		if text != "" {
			return text, true
		}
	}

	w.complete = true
	_ = w.f.Close()
	return "", false
}

func (w *wordlistIterator) Count() int {
	return -1
}

func (w *wordlistIterator) Complete() bool {
	return w.complete
}

func (w *wordlistIterator) Type() Type {
	return NewType(TypeList, InternalTypeWordlist)
}

func (w *wordlistIterator) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.filename)
}

func (w *wordlistIterator) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &w.filename)
}

type csvIterator struct {
	items []string
	ptr   int
}

func NewCSVIterator(raw string) *csvIterator {
	return &csvIterator{
		items: strings.Split(raw, ","),
	}
}

func (c *csvIterator) Next() (string, bool) {
	if c.ptr < len(c.items) {
		item := c.items[c.ptr]
		c.ptr++
		return item, true
	}
	return "", false
}

func (c *csvIterator) Count() int {
	return len(c.items)
}

func (c *csvIterator) Complete() bool {
	return c.ptr >= len(c.items)
}

func (c *csvIterator) Type() Type {
	return NewType(TypeList, InternalTypeCommaSeparatedList)
}

func (c *csvIterator) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Join(c.items, ","))
}

func (c *csvIterator) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.items = strings.Split(raw, ",")
	return nil
}
