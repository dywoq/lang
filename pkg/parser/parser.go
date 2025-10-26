package parser

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dywoq/lang/pkg/ast"
	"github.com/dywoq/lang/pkg/token"
)

type mini func(tree *ast.Tree) (ast.Node, error)

type Parser struct {
	tokens             []*token.Token
	parsing, lazysetup bool
	minis              []mini
	pos                int
	d                  debug
}

type debug struct {
	p  *Parser
	w  io.Writer
	on bool
}

// New returns a new pointer to Parser,
// with the automatically turned off debugger.
func New(tokens []*token.Token) *Parser {
	p := &Parser{
		tokens:  tokens,
		parsing: false,
		minis:   []mini{},
		pos:     0,
	}
	p.d = debug{p, nil, false}
	return p
}

func NewDebug(tokens []*token.Token, w io.Writer) *Parser {
	p := &Parser{
		tokens:  tokens,
		parsing: false,
		minis:   []mini{},
		pos:     0,
	}
	p.d = debug{p, w, true}
	return p
}

var ErrWorking = errors.New("parser: parsing right now")

var errEof = errors.New("parser: internal: eof")
var errOutOfBounds = errors.New("parser: internal: out of bounds")

func (p *Parser) Update(tokens []*token.Token) error {
	if p.parsing {
		return ErrWorking
	}
	if tokens == nil {
		return errors.New("parser: given tokens slice is nil")
	}
	for _, token := range tokens {
		if token == nil {
			return errors.New("parser: detected nil token when updating")
		}
	}
	p.tokens = tokens
	return nil
}

func (p *Parser) Parse() (*ast.Tree, error) {
	if len(p.tokens) == 0 {
		return nil, nil
	}
	p.debug("starting parsing")
	p.parsing = true
	defer func() {
		p.parsing = false
		p.debug("parsing ended")
	}()
	p.setup()
	p.reset()
	tree := &ast.Tree{}
	for {
		t, err := p.current()
		if err != nil {
			return nil, err
		}
		if t.Kind == token.Eof {
			break
		}
		n, err := p.parse(tree)
		if err != nil {
			return nil, err
		}
		tree.Statements = append(tree.Statements, n)
	}
	return tree, nil
}

func (p *Parser) setup() {
	if !p.lazysetup {
		p.minis = []mini{
			p.parseDeclaration,
		}
		p.lazysetup = true
	}
}

func (p *Parser) reset() {
	p.pos = 0
}

func (p *Parser) parse(tree *ast.Tree) (ast.Node, error) {
	chosen := ast.Node(nil)
	gotErr := error(nil)
	for _, mini := range p.minis {
		p.debug("parsing token")
		if chosen != nil {
			break
		}
		n, err := mini(tree)
		if err != nil {
			gotErr = err
			break
		}
		chosen = n
		p.debug("successfully parsed token")
	}
	return chosen, gotErr
}

func (p *Parser) current() (*token.Token, error) {
	if p.pos >= len(p.tokens) {
		return nil, errEof
	}
	p.debugf("getting current token: %s - %s", p.tokens[p.pos].Literal, p.tokens[p.pos].Kind)
	return p.tokens[p.pos], nil
}

func (p *Parser) exceptKind(k token.Kind) (*token.Token, error) {
	t, _ := p.current()
	if k != t.Kind {
		return nil, p.errorf("expected \"%s\" kind", k)
	}
	p.advance(1)
	return t, nil
}

func (p *Parser) exceptLiteral(lit string) (*token.Token, error) {
	t, _ := p.current()
	if lit != t.Literal {
		return nil, p.errorf("expected \"%s\" literal", lit)
	}
	p.advance(1)
	return t, nil
}

func (p *Parser) errorf(format string, v ...any) error {
	t, _ := p.current()
	return fmt.Errorf("%s; source is %d:%d", fmt.Sprintf(format, v...), t.Position.Line, t.Position.Column)
}

func (p *Parser) advance(n int) error {
	if p.pos+n >= len(p.tokens) {
		return errOutOfBounds
	}
	p.pos += n
	return nil
}

func (p *Parser) debug(v ...any) error {
	if !p.d.on {
		return nil
	}
	_, err := io.WriteString(p.d.w, fmt.Sprintf("%s %v\n", time.Now().String(), v))
	if err != nil {
		return err
	}
	return err
}

func (p *Parser) debugf(format string, v ...any) error {
	res := fmt.Sprintf(format, v...)
	return p.debug(res)
}

// mini parsers

func (p *Parser) parseValue() (ast.Node, error) {
	t, err := p.current()
	if err != nil {
		return nil, err
	}
	switch t.Kind {
	case token.Integer, token.Float, token.Identifier:
		return ast.Value{Value: t.Literal}, nil
	}
	return nil, p.errorf("unknown token kind: \"%s\"", t.Kind)
}

func (p *Parser) parseDeclaration(a *ast.Tree) (ast.Node, error) {
	ident, err := p.exceptKind(token.Identifier)
	if err != nil {
		return nil, err
	}
	tType, err := p.exceptKind(token.Type)
	if err != nil {
		return nil, err
	}
	val, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	a.Global = append(a.Global, ident.Literal)
	p.advance(1)
	return ast.Variable{
		Identifier: ident.Literal,
		Exported:   false,
		Const:      false,
		Consteval:  false,
		Copied:     false,
		Type:       tType.Literal,
		Value:      val,
	}, nil
}

func (d *debug) Set(b bool) error {
	if d.p.parsing {
		return ErrWorking
	}
	d.p.parsing = b
	return nil
}

func (d *debug) On() bool {
	return d.on
}

func (d *debug) SetWriter(w io.Writer) error {
	if d.p.parsing {
		return ErrWorking
	}
	if w == nil {
		return errors.New("debug: given io.Writer is nil")
	}
	d.w = w
	return nil
}
