package highlight

import "github.com/alecthomas/chroma"

var lightTheme = chroma.MustNewStyle("light", chroma.StyleEntries{
	chroma.Other:        "#ffffff",
	chroma.Background:   "#1d2432",
	chroma.Keyword:      "#ff636f",
	chroma.Name:         "#58a1dd",
	chroma.Literal:      "#a6be9d",
	chroma.Operator:     "#ff636f",
	chroma.OperatorWord: "#ff636f",
	chroma.Comment:      "italic #828b96",
})
