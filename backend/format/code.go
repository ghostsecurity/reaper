package format

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"

	"github.com/yosssi/gohtml"
)

func Code(input string, contentType string) string {
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	if contentType == "" {
		contentType = sniff(input)
	}
	switch contentType {
	case "":
		return input
	case "text/html":
		return gohtml.Format(input)
	case "application/xml":
		data, err := formatXML([]byte(input))
		if err != nil {
			return input
		}
		return string(data)
	case "application/json":
		var buffer bytes.Buffer
		if err := json.Indent(&buffer, []byte(input), "", "  "); err != nil {
			return input
		}
		return buffer.String()
	default:
		return input
	}
}

func sniff(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	switch {
	case strings.HasPrefix(input, "{"), strings.HasPrefix(input, "["):
		return "application/json"
	case strings.HasPrefix(input, "<!DOC"), strings.HasPrefix(input, "<html"):
		return "text/html"
	default:
		return ""
	}
}

func formatXML(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(b)
	encoder.Indent("", "  ")
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			if err := encoder.Flush(); err != nil {
				return nil, err
			}
			return b.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, err
		}
	}
}
