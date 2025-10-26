package token

import (
	"slices"
	"unicode"
)

// Kind represents the token kind.
type Kind string

// Collection represents the collection of tokens.
// To check if the string is present in collection,
// you can use slices.Contains:
//
//	c := token.Collection{"var", "let", "token"}
//	str := "var"
//	if !slices.Contains(c, str) {
//		panic("var is not in the collection")
//	}
type Collection []string

// Token is a structure with the
// literal, the given position and kind.
type Token struct {
	Literal  string    `json:"literal"`
	Kind     Kind      `json:"kind"`
	Position *Position `json:"position"`
}

// Position represents the token position.
type Position struct {
	Line     int `json:"line"`
	Column   int `json:"column"`
	Position int `json:"position"`
}

// A token kind.
const (
	Keyword            Kind = "keyword"
	ModifierConversion Kind = "modifier-conversion"
	Integer            Kind = "integer"
	Float              Kind = "float"
	String             Kind = "string"
	Separator          Kind = "separator"
	Illegal            Kind = "illegal"
	Type               Kind = "type"
	Identifier         Kind = "identifier"
	Eof                Kind = "eof"
)

var (
	Keywords = Collection{}

	Separators = Collection{
		",",
		";",
		"{",
		"}",
		"(",
		")",
		".",
	}

	ModifierConversions = Collection{
		"const",
		"consteval",
		"copy",
		"export",
	}

	Types = Collection{
		"i8",
		"i16",
		"i32",
		"i64",
		"u8",
		"u16",
		"u32",
		"u64",
		"str",
		"bool",
		"void",

		"uptr",
		"f32",
		"f64",

		"i128",
		"u128",

		"fix64",

		"f32",
		"f64",
	}
)

// New returns a pointer to Token.
func New(literal string, kind Kind, position *Position) *Token {
	return &Token{literal, kind, position}
}

// NewPosition returns a pointer to Position.
func NewPosition(line, column, position int) *Position {
	return &Position{line, column, position}
}

// IsIdentifier checks whether str is a valid identifier:
//   - Must not start with digit;
//   - Must start with letter of an underscore;
//   - No spaces or special characters;
//   - Cannot be a keyword name;
//   - Limit of the identifier is 255 symbols.
func IsIdentifier(str string) bool {
	if len(str) > 255 || len(str) == 0 {
		return false
	}
	if slices.Contains(Keywords, str) {
		return false
	}
	for i, r := range str {
		if i == 0 && unicode.IsDigit(r) {
			return false
		}
		if unicode.IsSymbol(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Valid checks whether k is a valid token kind.
func (k Kind) Valid() bool {
	kinds := []Kind{Keyword, ModifierConversion, Integer, Float, String, Separator}
	return slices.Contains(kinds, k)
}
