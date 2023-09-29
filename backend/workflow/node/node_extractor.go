package node

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/context"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"

	"github.com/ghostsecurity/reaper/backend/packaging"
	"github.com/ghostsecurity/reaper/backend/workflow/transmission"
)

type ExtractorNode struct {
	*base
	noInjections
}

const (
	ExtractBody      = "body"
	ExtractXPathXML  = "xpath_xml"
	ExtractXPathHTML = "xpath_html"
	ExtractXPathJSON = "xpath_json"
	ExtractRegexBody = "regex_body"
	ExtractHeader    = "header"
	ExtractStatus    = "status"
	ExtractSession   = "session"
)

func NewExtractor() *ExtractorNode {
	return &ExtractorNode{
		base: newBase(
			"Extractor",
			TypeExtractor,
			false,
			NewVarStorage(
				Connectors{
					NewConnector("response", transmission.TypeResponse|transmission.TypeMap, true),
					NewConnector("type", transmission.TypeChoice, false, ""),
					NewConnector("pattern", transmission.TypeString, false, ""),
					NewConnector("variable", transmission.TypeString, false, ""),
					NewConnector("strip", transmission.TypeBoolean, false, "Remove leading/trailing whitespace"),
				},
				Connectors{
					NewConnector("output", transmission.TypeMap, true),
				},
				map[string]transmission.Transmission{
					"type": transmission.NewChoice("body", map[string]string{
						ExtractBody:      "Entire Body",
						ExtractXPathHTML: "XPath: HTML",
						ExtractXPathXML:  "XPath: XML",
						ExtractXPathJSON: "XPath: JSON",
						ExtractRegexBody: "Body Regular Expression",
						ExtractHeader:    "Header",
						ExtractStatus:    "Status Code",
						ExtractSession:   "Session",
					}),
					"pattern":  transmission.NewString(""),
					"variable": transmission.NewString("$EXTRACTED$"),
					"strip":    transmission.NewBoolean(true),
				},
			),
		),
	}
}

func (n *ExtractorNode) Start(ctx context.Context, in <-chan Input, out chan<- OutputInstance, _ chan<- Output) error {
	strip, err := n.ReadInputBool("strip", nil)
	if err != nil {
		return err
	}

	extractType, err := n.ReadInputChoice("type", nil)
	if err != nil {
		return err
	}

	pattern, err := n.ReadInputString("pattern", nil)
	if err != nil {
		return err
	}

	variable, err := n.ReadInputString("variable", nil)
	if err != nil {
		return err
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

			response, err := n.ReadInputResponse("response", input.Data)
			if err != nil {
				return err
			}

			vars, err := n.ReadInputMap("response", input.Data)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			value, err := n.extract(response, extractType, pattern)
			if err != nil {
				return err
			}

			if strip {
				value = strings.TrimSpace(value)
			}

			if vars == nil {
				vars = make(map[string]string)
			}
			vars[variable] = value

			n.tryOut(ctx, out, OutputInstance{
				OutputName: "output",
				Complete:   input.Last,
				Data:       transmission.NewMap(vars),
			})
		}
	}
}

func (n *ExtractorNode) extract(response *packaging.HttpResponse, eType string, pattern string) (string, error) {
	switch eType {
	case ExtractBody:
		return response.Body, nil
	case ExtractStatus:
		return fmt.Sprintf("%d", response.StatusCode), nil
	case ExtractHeader:
		for _, header := range response.Headers {
			if strings.EqualFold(header.Key, pattern) {
				return header.Value, nil
			}
		}
		return "", nil
	case ExtractRegexBody:
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}
		matches := reg.FindAllStringSubmatch(response.Body, 1)
		if len(matches) == 0 {
			return "", nil
		}
		if len(matches[0]) == 0 {
			return "", nil
		}
		if len(matches[0]) == 1 {
			return matches[0][0], nil
		}
		return matches[0][1], nil
	case ExtractXPathHTML:
		return n.extractXPathHTML(response, pattern)
	case ExtractXPathXML:
		return n.extractXPathXML(response, pattern)
	case ExtractXPathJSON:
		return n.extractXPathJSON(response, pattern)
	case ExtractSession:
		var cookies []string
		for _, cookie := range response.Cookies {
			cookies = append(cookies, cookie.String())
		}
		return strings.Join(cookies, "; "), nil
	default:
		return "", fmt.Errorf("invalid extraction type: %s", eType)
	}
}

func (n *ExtractorNode) extractXPathHTML(response *packaging.HttpResponse, pattern string) (string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(response.Body))
	if err != nil {
		return "", err
	}

	node, err := htmlquery.Query(doc, pattern)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", nil
	}

	if node.FirstChild != nil {
		return node.FirstChild.Data, nil
	}

	return node.Data, nil
}

func (n *ExtractorNode) extractXPathXML(response *packaging.HttpResponse, pattern string) (string, error) {
	doc, err := xmlquery.Parse(strings.NewReader(response.Body))
	if err != nil {
		return "", err
	}

	node, err := xmlquery.Query(doc, pattern)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", nil
	}

	return node.InnerText(), nil
}

func (n *ExtractorNode) extractXPathJSON(response *packaging.HttpResponse, pattern string) (string, error) {
	doc, err := jsonquery.Parse(strings.NewReader(response.Body))
	if err != nil {
		return "", err
	}

	node, err := jsonquery.Query(doc, pattern)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", nil
	}

	return fmt.Sprintf("%s", node.Value()), nil
}
