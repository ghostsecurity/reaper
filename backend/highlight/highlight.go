package highlight

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
)

var customStyle = chroma.MustNewStyle("ghost", chroma.StyleEntries{
	chroma.Error:                 "#bf616a",
	chroma.Background:            "#d8dee9 bg:#2e3440",
	chroma.Keyword:               "bold #81a1c1",
	chroma.KeywordPseudo:         "nobold #81a1c1",
	chroma.KeywordType:           "nobold #81a1c1",
	chroma.Name:                  "#eceff4",
	chroma.NameAttribute:         "#8fbcbb",
	chroma.NameBuiltin:           "#81a1c1",
	chroma.NameClass:             "#8fbcbb",
	chroma.NameConstant:          "#8fbcbb",
	chroma.NameDecorator:         "#d08770",
	chroma.NameEntity:            "#d08770",
	chroma.NameException:         "#bf616a",
	chroma.NameFunction:          "#88c0d0",
	chroma.NameLabel:             "#8fbcbb",
	chroma.NameNamespace:         "#8fbcbb",
	chroma.NameTag:               "#81a1c1",
	chroma.NameVariable:          "#d8dee9",
	chroma.Literal:               "#d8dee9",
	chroma.LiteralString:         "#a3be8c",
	chroma.LiteralStringDoc:      "#616e87",
	chroma.LiteralStringEscape:   "#ebcb8b",
	chroma.LiteralStringInterpol: "#a3be8c",
	chroma.LiteralStringOther:    "#a3be8c",
	chroma.LiteralStringRegex:    "#ebcb8b",
	chroma.LiteralStringSymbol:   "#a3be8c",
	chroma.LiteralNumber:         "#b48ead",
	chroma.Operator:              "#81a1c1",
	chroma.OperatorWord:          "bold #81a1c1",
	chroma.Punctuation:           "#eceff4",
	chroma.Comment:               "italic #616e87",
	chroma.CommentPreproc:        "#5e81ac",
	chroma.GenericDeleted:        "#bf616a",
	chroma.GenericEmph:           "italic",
	chroma.GenericError:          "#bf616a",
	chroma.GenericHeading:        "bold #88c0d0",
	chroma.GenericInserted:       "#a3be8c",
	chroma.GenericOutput:         "#d8dee9",
	chroma.GenericPrompt:         "bold #4c566a",
	chroma.GenericStrong:         "bold",
	chroma.GenericSubheading:     "bold #88c0d0",
	chroma.GenericTraceback:      "#bf616a",
	chroma.TextWhitespace:        "#d8dee9",
})

func HTTP(raw string) string {

	raw = strings.ReplaceAll(raw, "\r\n", "\n")

	headers, body, _ := strings.Cut(raw, "\n\n")

	highlightedHeaders, err := highlightHeaders(headers, customStyle)
	if err != nil {
		return ""
	}

	var ct string
	for _, line := range strings.Split(headers, "\n") {
		line = strings.TrimSpace(strings.ToLower(line))
		if strings.HasPrefix(line, "content-type:") {
			ct = strings.TrimSpace(strings.TrimPrefix(line, "content-type:"))
			break
		}
	}

	buf := bytes.NewBuffer(nil)

	if len(headers) > 0 {
		buf.Write([]byte("<br/>"))
		if len(body) == 0 {
			buf.Write([]byte("<br/>"))
		}
	}

	return highlightedHeaders + buf.String() + Body(body, ct)
}

func Body(body string, contentType string) string {

	var bodyLexer chroma.Lexer
	switch {
	case strings.HasPrefix(contentType, "text/html"):
		bodyLexer = lexers.Get("html")
	case strings.HasPrefix(contentType, "text/javascript"):
		bodyLexer = lexers.Get("javascript")
	case strings.HasPrefix(contentType, "text/css"):
		bodyLexer = lexers.Get("css")
	case strings.HasPrefix(contentType, "text/xml"):
		bodyLexer = lexers.Get("xml")
	case strings.Contains(contentType, "json"):
		bodyLexer = lexers.Get("json")
	default:
		bodyLexer = lexers.Get("plaintext")
	}

	if bodyLexer == nil {
		bodyLexer = lexers.Fallback
	}

	bodyIter, err := bodyLexer.Tokenise(&chroma.TokeniseOptions{
		State:    "root",
		Nested:   false,
		EnsureLF: false,
	}, body)
	if err != nil {
		return ""
	}

	buf := bytes.NewBuffer(nil)

	formatter := html.New(html.PreventSurroundingPre(true))
	if err := formatter.Format(buf, customStyle, bodyIter); err != nil {
		return ""
	}

	// TODO: do we need an extra <br/> here when empty?
	return buf.String()
}

func highlightHeaders(raw string, style *chroma.Style) (string, error) {

	buf := bytes.NewBuffer(nil)

	httpLexer := lexers.Get("http")
	if httpLexer == nil {
		httpLexer = lexers.Fallback
	}
	formatter := html.New(html.PreventSurroundingPre(true))

	httpIter, err := httpLexer.Tokenise(&chroma.TokeniseOptions{
		State:    "root",
		Nested:   false,
		EnsureLF: false,
	}, raw)
	if err != nil {
		return "", fmt.Errorf("error tokenising headers: %s", err)
	}

	if err := formatter.Format(buf, style, httpIter); err != nil {
		return "", fmt.Errorf("error formatting headers: %s", err)
	}

	return buf.String(), nil
}
