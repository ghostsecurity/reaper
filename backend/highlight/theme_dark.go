package highlight

import "github.com/alecthomas/chroma"

var darkTheme = chroma.MustNewStyle("dark", chroma.StyleEntries{
	chroma.Error:             "#ff0000",
	chroma.Background:        "#e5e5e5 bg:#000000",
	chroma.Keyword:           "bold #ffffff",
	chroma.Name:              "#2299cf",
	chroma.NameAttribute:     "#007f7f",
	chroma.NameBuiltin:       "bold #ffffff",
	chroma.NameException:     "#22ddee",
	chroma.NameKeyword:       "bold #ffffff",
	chroma.NameTag:           "bold",
	chroma.LiteralDate:       "bold #22ddee",
	chroma.LiteralString:     "bold #00ffff",
	chroma.LiteralNumber:     "bold #22ddee",
	chroma.Comment:           "#007f7f",
	chroma.CommentPreproc:    "bold #00ff00",
	chroma.GenericHeading:    "bold",
	chroma.GenericStrong:     "bold",
	chroma.GenericSubheading: "bold",
	chroma.GenericUnderline:  "underline",
})
