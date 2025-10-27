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
	fileName           string
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

// NewDebug returns a new pointer to Parser,
// with the automatically turned on debugger.
//
// w is used by the debugger, to write the messages.
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

// ErrWorking means the parser is currently parsing,
// so the debug options can't be changed,
// and tokens can't be updated.
var ErrWorking = errors.New("parser: parsing right now")

var errEof = errors.New("parser: internal: eof")
var errOutOfBounds = errors.New("parser: internal: out of bounds")

// Update updates the underlying tokens to tokens.
// Returns ErrWorking if the parser is currently working,
// or an error tokens slice is nil.
//
// A tokens of tokens slice must be not nil,
// otherwise an error is returned.
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

// Parse parses the tokens into ast.Tree,
// returning a pointer to ast.Tree.
//
// Returns an error if the got tokens are empty,
// and any errors encountered during parsing.
func (p *Parser) Parse(fileName string) (*ast.Tree, error) {
	if len(p.tokens) == 0 {
		return nil, errors.New("parser: got no tokens")
	}
	p.fileName = fileName
	p.debug("starting parsing")
	p.parsing = true
	defer func() {
		p.parsing = false
		p.debug("parsing ended")
	}()
	p.setup()
	p.reset()
	tree := &ast.Tree{}
	tree.File = p.fileName
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

func (p *Parser) expectKind(k token.Kind) (*token.Token, error) {
	t, _ := p.current()
	if k != t.Kind {
		return nil, p.errorf("expected \"%s\" kind, got \"%s\" (literal: %s) instead", k, t.Kind, t.Literal)
	}
	p.advance(1)
	return t, nil
}

func (p *Parser) expectLiteral(lit string) (*token.Token, error) {
	t, _ := p.current()
	if lit != t.Literal {
		return nil, p.errorf("expected \"%s\" literal, got \"%s\" instead", lit, t.Literal)
	}
	p.advance(1)
	return t, nil
}

func (p *Parser) errorf(format string, v ...any) error {
	t, _ := p.current()
	return fmt.Errorf("%s; source is %s:%d", fmt.Sprintf(format, v...), p.fileName, t.Position.Line)
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
	case token.Integer, token.Float, token.Identifier, token.String:
		p.advance(1)
		return ast.Value{Value: t.Literal, Kind: t.Kind}, nil

	case token.ModifierConversion:
		name, _ := p.current()
		p.advance(1)
		_, err = p.expectLiteral("(")
		if err != nil {
			return nil, err
		}
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		_, err = p.expectLiteral(")")
		if err != nil {
			return nil, err
		}
		return ast.ModifierConversion{Name: name.Literal, Value: val}, nil
	}

	if t.Literal == "(" {
		return p.parseFunctionValue()
	}

	return nil, p.errorf("unknown token kind: \"%s\"", t.Kind)
}

func (p *Parser) parseFunctionValue() (ast.Node, error) {
	_, err := p.expectLiteral("(")
	if err != nil {
		return nil, err
	}
	args := []ast.FunctionArgument{}
	for {
		t, _ := p.current()
		if t.Literal == ")" {
			break
		}
		if t.Literal == "," {
			p.advance(1)
			continue
		}
		variadic := false
		name, err := p.expectKind(token.Identifier)
		if err != nil {
			return nil, err
		}
		tType, err := p.expectKind(token.Type)
		if err != nil {
			return nil, err
		}

		if t, _ := p.current(); t.Literal == "^" {
			variadic = true
			p.advance(1)
		}

		args = append(args, ast.FunctionArgument{
			Name:     name.Literal,
			Type:     tType.Literal,
			Variadic: variadic,
		})
	}
	_, err = p.expectLiteral(")")
	if err != nil {
		return nil, err
	}

	body, err := p.parseFunctionBody()
	if err != nil {
		return nil, err
	}

	return ast.Function{
		Args: args,
		Body: body,
	}, nil
}

func (p *Parser) parseFunctionBody() ([]ast.Node, error) {
	_, err := p.expectLiteral("{")
	if err != nil {
		return nil, err
	}
	body := []ast.Node{}
	for {
		t, _ := p.current()
		if t.Literal == "}" {
			p.advance(1)
			break
		}
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, statement)
	}
	return body, nil
}

func (p *Parser) parseStatement() (ast.Node, error) {
	name, err := p.expectKind(token.Identifier)
	if err != nil {
		return nil, err
	}
	args := []ast.InstructionArgument{}
	for {
		t, _ := p.current()
		if t.Literal == "," {
			p.advance(1)
			continue
		}
		if t.Literal == ";" {
			p.advance(1)
			break
		}
		arg, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		args = append(args, ast.InstructionArgument{Value: arg})
	}
	return ast.Instruction{
		Name: name.Literal,
		Args: args,
	}, nil
}

func (p *Parser) parseDeclaration(a *ast.Tree) (ast.Node, error) {
	ident, err := p.expectKind(token.Identifier)
	if err != nil {
		return nil, err
	}
	tType, err := p.expectKind(token.Type)
	if err != nil {
		return nil, err
	}
	val, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	a.Global = append(a.Global, ident.Literal)
	return ast.Declaration{
		Identifier: ident.Literal,
		Type:       tType.Literal,
		Value:      val,
	}, nil
}

// Set turns on the debugging mode.
// Returns ErrWorking if the parser is working right now.
func (d *debug) Set(b bool) error {
	if d.p.parsing {
		return ErrWorking
	}
	d.p.parsing = b
	return nil
}

// On returns true if the debugging mode is on.
func (d *debug) On() bool {
	return d.on
}

// SetWriter sets a instance that implements io.Writer interface.
// Returns ErrWorking if the parser is working right now.
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
