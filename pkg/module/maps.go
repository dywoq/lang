package module

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// MapParser is responsible for parsing modules map.
type MapParser struct {
	r     io.Reader
	input []byte
	pos   int
}

type MapSymbol struct {
	FullPath string `json:"full_path"`
	Chain    *Chain `json:"chain"`
	Module   string `json:"module"`
}

// ErrUnexpectedEof is returned by module maps parser,
// when it was parsing and met EOF (End Of File).
var ErrUnexpectedEof = errors.New("module: unexpected eof")

// Parse parses the bytes got from r with io.ReadAll.
//
// Returns an error if it encountered an unexpected EOF,
// or the parser didn't get what expected.
//
// Returns a map of symbols, where the key is the name of symbol,
// and
func (m *MapParser) Parse(r io.Reader) (map[string]*MapSymbol, error) {
	m.r = r
	err := m.update()
	if err != nil {
		return nil, err
	}
	if len(m.input) == 0 {
		return nil, nil
	}
	symbols := map[string]*MapSymbol{}
	for !m.eof() {
		symbol, err := m.parse()
		if err != nil {
			return nil, err
		}
		symbols[m.parseName(symbol)] = symbol
	}
	return symbols, nil
}

func (m *MapParser) update() error {
	bytes, err := io.ReadAll(m.r)
	if err != nil {
		return err
	}
	m.input = bytes
	return nil
}

func (m *MapParser) eof() bool {
	return m.pos >= len(m.input)
}

func (m *MapParser) parse() (*MapSymbol, error) {
	s := &MapSymbol{Chain: &Chain{}}
	ptr := s.Chain
	for !m.eof() && unicode.IsSpace(rune(m.input[m.pos])) {
		m.pos++
	}
	var nameStart int
	for {
		nameStart = m.pos
		for !m.eof() && (unicode.IsLetter(rune(m.input[m.pos])) || unicode.IsDigit(rune(m.input[m.pos])) || m.input[m.pos] == '_') {
			m.pos++
		}
		if nameStart == m.pos {
			return nil, fmt.Errorf("module: expected identifier, got %q", m.input[m.pos])
		}
		ptr.Name = string(m.input[nameStart:m.pos])
		if m.eof() {
			return nil, ErrUnexpectedEof
		}
		if m.input[m.pos] == '.' {
			m.pos++
			ptr.Next = &Chain{}
			ptr = ptr.Next
			continue
		}
		break
	}
	for !m.eof() && unicode.IsSpace(rune(m.input[m.pos])) {
		m.pos++
	}
	if m.eof() || m.input[m.pos] != ':' {
		return nil, fmt.Errorf("module: expected ':' after identifier, got %q", m.input[m.pos])
	}
	m.pos++
	for !m.eof() && unicode.IsSpace(rune(m.input[m.pos])) {
		m.pos++
	}
	if m.eof() || m.input[m.pos] != '"' {
		return nil, errors.New(`module: expected '"' after ':'`)
	}
	m.pos++
	pathStart := m.pos
	for !m.eof() && m.input[m.pos] != '"' {
		m.pos++
	}
	if m.eof() {
		return nil, ErrUnexpectedEof
	}
	fullPath := string(m.input[pathStart:m.pos])
	m.pos++
	path, err := m.parsePath(fullPath)
	if err != nil {
		return nil, err
	}
	s.FullPath = path
	s.Module = s.Chain.Name
	return s, nil
}

func (m *MapParser) parsePath(fullPath string) (string, error) {
	var b strings.Builder
	pos := 0
	for pos < len(fullPath) {
		r := rune(fullPath[pos])
		if r == '$' {
			pos++
			if pos >= len(fullPath) {
				return "", ErrUnexpectedEof
			}
			start := pos
			for pos < len(fullPath) && rune(fullPath[pos]) != '$' {
				pos++
			}
			if pos >= len(fullPath) {
				return "", errors.New("module: unterminated env variable ($... missing closing $)")
			}
			envName := fullPath[start:pos]
			pos++
			value := os.Getenv(envName)
			b.WriteString(value)
			continue
		}
		b.WriteRune(r)
		pos++
	}

	return b.String(), nil
}

func (m *MapParser) parseName(s *MapSymbol) string {
	var b strings.Builder
	for ptr := s.Chain; ptr != nil; ptr = ptr.Next {
		if b.Len() > 0 {
			b.WriteString(".")
		}
		b.WriteString(ptr.Name)
	}
	return b.String()
}
