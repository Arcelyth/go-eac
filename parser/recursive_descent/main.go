package main

import (
	"bytes"
	"fmt"
)

type Parser struct {
	buffer *bytes.Buffer // use bytes instead of <lexeme, ty> in order to simplify
	cur    byte
}

func newParser(s string) Parser {
	return Parser{
		buffer: bytes.NewBufferString(s),
		cur:    0,
	}
}

func (p *Parser) Next() bool {
	c := p.buffer.Next(1)
	if len(c) == 0 {
		p.cur = 0
		return false
	}
	p.cur = c[0]
	return true
}

func (p *Parser) TermPrime() bool {
	// Term' -> x Factor Term'
	// Term' -> / Factor Term'
	switch p.cur {
	case 'x', '/':
		p.Next()
		if p.Factor() {
			return p.TermPrime()
		} else {
			return false
		}
	// Term' -> null
	case '+', '-', ')', 0:
		return true
	default:
		return false
	}
}

func (p *Parser) Factor() bool {
	// Factor -> ( Expr )
	switch {
	case p.cur == '(':
		p.Next()
		if !p.Expr() {
			return false
		}
		if p.cur != ')' {
			return false
		}
		p.Next()
		return true
	// Factor -> num
	case p.cur >= '0' && p.cur <= '9': // only use one digit in order to simplify
		p.Next()
		return true
	default:
		return false
	}
}

func (p *Parser) Term() bool {
	// Term -> Factor Term'
	if p.Factor() {
		return p.TermPrime()
	} else {
		return false
	}
}

func (p *Parser) ExprPrime() bool {
	// Expr' -> + Term Expr'
	// Expr' -> - Term Expr'
	switch p.cur {
	case '+', '-':
		p.Next()
		if p.Term() {
			return p.ExprPrime()
		} else {
			return false
		}
	// Expr' -> null
	case ')', 0:
		return true
	default:
		return false
	}
}

func (p *Parser) Expr() bool {
	// Expr -> Term Expr'
	if p.Term() {
		return p.ExprPrime()
	} else {
		return false
	}
}

func (p *Parser) Failed() {
	fmt.Println("Failed")
}

func main() {
	src := "1+4x5"
	parser := newParser(src)
	ok := parser.Next()
	if parser.Expr() {
		if ok {
			fmt.Println("Success")
		} else {
			parser.Failed()
		}
	} else {
		parser.Failed()
	}

}
