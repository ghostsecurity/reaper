package highlight

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"strings"
)

func HTTP(raw string, themeName string) string {

	raw = strings.ReplaceAll(raw, "\r\n", "\n")

	headers, body, _ := strings.Cut(raw, "\n\n")

	style := findTheme(themeName)

	highlightedHeaders, err := highlightHeaders(headers, style)
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

	var bodyLexer chroma.Lexer
	switch {
	case strings.HasPrefix(ct, "text/html"):
		bodyLexer = lexers.Get("html")
	case strings.HasPrefix(ct, "text/javascript"):
		bodyLexer = lexers.Get("javascript")
	case strings.HasPrefix(ct, "text/css"):
		bodyLexer = lexers.Get("css")
	case strings.HasPrefix(ct, "text/xml"):
		bodyLexer = lexers.Get("xml")
	case strings.Contains(ct, "json"):
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

	if len(headers) > 0 {
		buf.Write([]byte("<br/>"))
		if len(body) == 0 {
			buf.Write([]byte("<br/>"))
		}
	}

	formatter := html.New(html.PreventSurroundingPre(true))
	if err := formatter.Format(buf, style, bodyIter); err != nil {
		return ""
	}
	return highlightedHeaders + buf.String()
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
