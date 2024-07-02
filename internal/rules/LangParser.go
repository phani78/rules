package rules

import (
	"errors"
	"fmt"
	"strconv"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenNumber
	tokenIdentifier
	tokenPlus
	tokenMinus
	tokenAsterisk
	tokenSlash
	tokenLParen
	tokenRParen
	tokenLEQ // <=
	tokenGEQ // >=
	tokenLT  // <
	tokenGT  // >
	tokenEQ  // ==
	tokenNEQ // !=
)

type token struct {
	typ tokenType
	val string
}

type lexer struct {
	input string
	start int
	pos   int
	width int
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return rune(tokenEOF)
	}
	r := rune(l.input[l.pos])
	l.width = 1
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(t tokenType) token {
	tok := token{t, l.input[l.start:l.pos]}
	l.start = l.pos
	return tok
}

func (l *lexer) skipWhitespace() {
	for {
		r := l.next()
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			l.backup()
			break
		}
		l.start = l.pos
	}
}

func (l *lexer) lex() token {
	l.skipWhitespace()
	switch r := l.next(); {
	case r == rune(tokenEOF):
		return l.emit(tokenEOF)
	case r == '+':
		return l.emit(tokenPlus)
	case r == '-':
		return l.emit(tokenMinus)
	case r == '*':
		return l.emit(tokenAsterisk)
	case r == '/':
		return l.emit(tokenSlash)
	case r == '(':
		return l.emit(tokenLParen)
	case r == ')':
		return l.emit(tokenRParen)
	case r == '<':
		if l.next() == '=' {
			return l.emit(tokenLEQ)
		}
		l.backup()
		return l.emit(tokenLT)
	case r == '>':
		if l.next() == '=' {
			return l.emit(tokenGEQ)
		}
		l.backup()
		return l.emit(tokenGT)
	case r == '=':
		if l.next() == '=' {
			return l.emit(tokenEQ)
		}
		l.backup()
		return l.emit(tokenError)
	case r == '!':
		if l.next() == '=' {
			return l.emit(tokenNEQ)
		}
		l.backup()
		return l.emit(tokenError)
	case '0' <= r && r <= '9':
		l.backup()
		return l.lexNumber()
	case 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_':
		l.backup()
		return l.lexIdentifier()
	default:
		return l.emit(tokenError)
	}
}

func (l *lexer) lexNumber() token {
	for {
		if r := l.next(); r < '0' || r > '9' {
			l.backup()
			break
		}
	}
	return l.emit(tokenNumber)
}

func (l *lexer) lexIdentifier() token {
	for {
		if r := l.next(); !('a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '_') {
			l.backup()
			break
		}
	}
	return l.emit(tokenIdentifier)
}

type parser struct {
	tokens []token
	pos    int
	vars   map[string]int
}

func (p *parser) next() token {
	if p.pos >= len(p.tokens) {
		return token{tokenEOF, ""}
	}
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *parser) backup() {
	p.pos--
}

func (p *parser) parsePrimary() (int, error) {
	tok := p.next()
	switch tok.typ {
	case tokenNumber:
		return strconv.Atoi(tok.val)
	case tokenIdentifier:
		val, ok := p.vars[tok.val]
		if !ok {
			return 0, fmt.Errorf("undefined variable: %s", tok.val)
		}
		return val, nil
	case tokenLParen:
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		if tok := p.next(); tok.typ != tokenRParen {
			return 0, errors.New("expected closing parenthesis")
		}
		return val, nil
	default:
		return 0, fmt.Errorf("unexpected token: %s", tok.val)
	}
}

func (p *parser) parseUnary() (int, error) {
	tok := p.next()
	if tok.typ == tokenPlus {
		return p.parsePrimary()
	} else if tok.typ == tokenMinus {
		val, err := p.parsePrimary()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	p.backup()
	return p.parsePrimary()
}

func (p *parser) parseMultiplicative() (int, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	for {
		tok := p.next()
		switch tok.typ {
		case tokenAsterisk:
			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			left *= right
		case tokenSlash:
			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			left /= right
		default:
			p.backup()
			return left, nil
		}
	}
}

func (p *parser) parseAdditive() (int, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return 0, err
	}
	for {
		tok := p.next()
		switch tok.typ {
		case tokenPlus:
			right, err := p.parseMultiplicative()
			if err != nil {
				return 0, err
			}
			left += right
		case tokenMinus:
			right, err := p.parseMultiplicative()
			if err != nil {
				return 0, err
			}
			left -= right
		default:
			p.backup()
			return left, nil
		}
	}
}

func (p *parser) parseComparison() (bool, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return false, err
	}
	tok := p.next()
	switch tok.typ {
	case tokenLEQ:
		right, err := p.parseAdditive()
		if err != nil {
			return false, err
		}
		return left <= right, nil
	case tokenGEQ:
		right, err := p.parseAdditive()
		if err != nil {
			return false, err
		}
		return left >= right, nil
	case tokenLT:
		right, err := p.parseAdditive()
		if err != nil {
			return false, err
		}
		return left < right, nil
	case tokenGT:
		right, err := p.parseAdditive()
		if err != nil {
			return false, err
		}
		return left > right, nil
	case tokenEQ:
		right, err := p.parseAdditive()
		if err != nil {
			return false, err
		}
		return left == right, nil
	case tokenNEQ:
		right, err := p.parseAdditive()
		if err != nil {
			return false, err
		}
		return left != right, nil
	default:
		return false, fmt.Errorf("unexpected token in comparison: %s", tok.val)
	}
}

func (p *parser) parseExpression() (int, error) {
	return p.parseAdditive()
}

func parseAndEvaluate(expr string, vars map[string]int) (bool, error) {
	l := &lexer{input: expr}
	tokens := []token{}
	for {
		tok := l.lex()
		if tok.typ == tokenEOF {
			break
		}
		if tok.typ == tokenError {
			return false, fmt.Errorf("lexing error: %s", tok.val)
		}
		tokens = append(tokens, tok)
	}
	p := &parser{tokens: tokens, vars: vars}
	return p.parseComparison()
}

func main() {
	vars := map[string]int{
		"currentSum": 10,
		"salesTotal": 4,
	}
	expr := "currentSum <= salesTotal + 5"
	result, err := parseAndEvaluate(expr, vars)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
