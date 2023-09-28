package node

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
)

type IfNode struct {
	*base
	noInjections
}

const (
	IfEqual          = "equal"
	IfNotEqual       = "equal_not"
	IfGreater        = "greater"
	IfGreaterOrEqual = "greater_or_equal"
	IfLess           = "less"
	IfLessOrEqual    = "less_or_equal"
	IfContains       = "contains"
	IfNotContains    = "contains_not"
	IfRegexMatch     = "regex_match"
	IfRegexNotMatch  = "regex_match_not"
)

func NewIf() *IfNode {
	return &IfNode{
		base: newBase(
			"If",
			TypeIf,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("vars", transmission.TypeMap, true),
					NewConnector("a", transmission.TypeString, false, ""),
					NewConnector("comparison", transmission.TypeChoice, false, ""),
					NewConnector("z", transmission.TypeString, false, ""),
				},
				Connectors{
					NewConnector("true", transmission.TypeMap, true),
					NewConnector("false", transmission.TypeMap, true),
				},
				map[string]transmission.Transmission{
					"comparison": transmission.NewChoice("comparison", map[string]string{
						IfEqual:          "==",
						IfNotEqual:       "!=",
						IfGreater:        ">",
						IfLess:           "<",
						IfGreaterOrEqual: ">=",
						IfLessOrEqual:    "<=",
						IfContains:       "contains",
						IfNotContains:    "not contains",
						IfRegexMatch:     "regex match",
						IfRegexNotMatch:  "regex no match",
					}),
					"a": transmission.NewString("$A$"),
					"z": transmission.NewString("$Z$"),
				},
			),
		),
	}
}

func (n *IfNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {

	comparison, err := n.ReadInputChoice("comparison", nil)
	if err != nil {
		return fmt.Errorf("invalid input 'comparison': %s", err)
	}

	defer n.setBusy(false)

	for {
		n.setBusy(false)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case input, ok := <-in:
			if !ok {
				return nil
			}

			n.setBusy(true)

			if input.Data == nil {
				return fmt.Errorf("input is nil")
			}

			vars, err := n.ReadInputMap("vars", input.Data)
			if err != nil {
				return fmt.Errorf("invalid input 'vars': %s", err)
			}
			if vars == nil {
				vars = make(map[string]string)
			}

			a, err := n.ReadInputString("a", input.Data)
			if err != nil {
				return fmt.Errorf("invalid input 'a': %s", err)
			}
			z, err := n.ReadInputString("z", input.Data)
			if err != nil {

			}

			for k, v := range vars {
				a = strings.ReplaceAll(a, k, v)
				z = strings.ReplaceAll(z, k, v)
			}

			result, err := n.compare(comparison, a, z)

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if result {
				n.tryOut(ctx, out, OutputInstance{
					OutputName: "true",
					Complete:   input.Last,
					Data:       transmission.NewMap(vars),
				})
			} else {
				n.tryOut(ctx, out, OutputInstance{
					OutputName: "false",
					Complete:   input.Last,
					Data:       transmission.NewMap(vars),
				})
			}
		}
	}
}

func (n *IfNode) compare(comparison string, a, z string) (bool, error) {
	switch comparison {
	case IfEqual:
		return a == z, nil
	case IfNotEqual:
		return a != z, nil
	case IfGreater:
		af, zf := n.numParams(a, z)
		return af > zf, nil
	case IfLess:
		af, zf := n.numParams(a, z)
		return af < zf, nil
	case IfGreaterOrEqual:
		af, zf := n.numParams(a, z)
		return af >= zf, nil
	case IfLessOrEqual:
		af, zf := n.numParams(a, z)
		return af <= zf, nil
	case IfContains:
		return strings.Contains(a, z), nil
	case IfNotContains:
		return !strings.Contains(a, z), nil
	case IfRegexMatch:
		return regexpMatch(a, z)
	case IfRegexNotMatch:
		m, err := regexpMatch(a, z)
		return !m, err
	default:
		return false, fmt.Errorf("invalid comparison: %s", comparison)
	}
}

func regexpMatch(a string, z string) (bool, error) {
	r, err := regexp.Compile(z)
	if err != nil {
		return false, err
	}
	return r.MatchString(a), nil
}

func (n *IfNode) numParams(a, z string) (float64, float64) {
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		af = 0
	}
	zf, err := strconv.ParseFloat(z, 64)
	if err != nil {
		zf = 0
	}
	return af, zf
}
