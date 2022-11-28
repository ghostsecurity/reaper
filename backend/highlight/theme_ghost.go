package highlight

import "github.com/alecthomas/chroma"

/*

From Ghost Branding Guidelines

main colours

white: #fff
charcoal:      #222531

secondary colours

steel grey:    #6B7292
cloud grey:    #C3C8DF
light grey:    #F2F4F7
midnight blue: #492CFB
lapis:         #4A64FA
papaya: #FA7400
carrot: #FF8922

data visualisation

red:           #ED475B
green:         #93D65E
orange:        #F5A42A
blue:          #58BAC8

custom extras:

6A84FA - lighter purple for dark bg

*/

var ghostTheme = chroma.MustNewStyle("ghost", chroma.StyleEntries{
	chroma.Error:             "#dd1111 bg:#1e0010", // Error
	chroma.Background:        "bg:#222531",         // Background (charcoal)
	chroma.Keyword:           "#6A84FA",            // Protocol name (e.g. HTTP part of HTTP/1.1)
	chroma.KeywordNamespace:  "#6A84FA",            // ?
	chroma.Name:              "#6A84FA",            // HTTP header keys
	chroma.NameFunction:      "bold #6A84FA",       // HTTP method
	chroma.NameNamespace:     "#C3C8DF",            // URI - e.g. in GET /xyz HTTP/1.1 (cloud grey)
	chroma.NameAttribute:     "#6A84FA",            // HTML attribute names
	chroma.NameBuiltinPseudo: "#6A84FA",
	chroma.NameClass:         "#6A84FA",
	chroma.NameConstant:      "#6A84FA",
	chroma.NameDecorator:     "#6A84FA",
	chroma.NameException:     "#6A84FA", // "OK" in HTTP/1.1 200 OK
	chroma.NameProperty:      "#6A84FA",
	chroma.NameTag:           "#6A84FA", // html tag names
	chroma.NameVariable:      "#6A84FA",
	chroma.Literal:           "#ffffff", // HTTP header values
	chroma.LiteralDate:       "#ffffff",
	chroma.LiteralString:     "#ffffff",
	chroma.LiteralNumber:     "#ffffff",
	chroma.Operator:          "#C3C8DF", // operators
	chroma.Punctuation:       "#C3C8DF", // punctuation
	chroma.Comment:           "#C3C8DF", // Comments (steel grey)
	chroma.GenericDeleted:    "#ED475B",
	chroma.GenericEmph:       "italic",
	chroma.GenericInserted:   "#93D65E",
	chroma.GenericStrong:     "bold",
	chroma.GenericSubheading: "#ffffff",
	chroma.Text:              "#ffffff", // default
	chroma.TextWhitespace:    "#C3C8DF",
})
