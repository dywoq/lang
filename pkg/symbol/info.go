package symbol

import (
	"errors"

	"github.com/dywoq/lang/pkg/ast"
	"github.com/dywoq/lang/pkg/token"
)

// Info is responsible for collecting the information
// of global symbols.
type Info struct {
	symbols map[string]*Symbol
}

// NewInfo returns a new pointer to Info.
func NewInfo() *Info {
	return &Info{make(map[string]*Symbol)}
}

// Collect collects the information about the global symbols,
// returning a map of symbols.
//
// All statements in t.Statements must be global symbols,
// otherwise the error is returned.
func (i *Info) Collect(t *ast.Tree) (map[string]*Symbol, error) {
	for key := range i.symbols {
		delete(i.symbols, key)
	}
	for _, statement := range t.Statements {
		s := &Symbol{}
		d, ok := statement.(ast.Declaration)
		if !ok {
			return nil, errors.New("symbol: info: expected declaration")
		}
		err := i.detectModifiers(s, d.Value)
		if err != nil {
			return nil, err
		}
		i.symbols[d.Identifier] = s
	}
	return i.symbols, nil
}

func (i *Info) detectModifiers(s *Symbol, n ast.Node) error {
	t, ok := n.(ast.ModifierConversion)
	if !ok {
		return nil
	}
	switch t.Name {
	case "export":
		s.Exported = true
	case "consteval":
		s.Consteval = true
	case "copy":
		s.Copied = true
		source, ok := t.Value.(ast.Value)
		if !ok || !token.IsIdentifier(source.Value) {
			return errors.New("symbol: info: expected identifier in copy(...)")
		}
		s.CopiedFrom = source.Value
	case "const":
		s.Const = true
	}
	if inner, ok := t.Value.(ast.ModifierConversion); ok {
		return i.detectModifiers(s, inner)
	}
	return nil
}
