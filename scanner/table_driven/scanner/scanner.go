package main

import (
	"fmt"
	"slices"
)

type charType byte
type State byte

const (
	S0 State = iota
	S1
	S2
	Se
	Bad
)

const (
	Register charType = iota
	Data
	Other
)

type Scanner struct {
	input string
	pos   int
	stack []State
}

func new_scanner(input string) Scanner {
	return Scanner{
		input: input,
		pos:   0,
	}
}

func (s *Scanner) CharCat(ch byte) charType {
	switch {
	case ch == 'r':
		return Register
	case ch >= '0' && ch <= '9':
		return Data
	default:
		return Other
	}
}

func (s *Scanner) NextChar() (byte, bool) {
	if s.pos >= len(s.input) {
		return 0, false
	}
	ch := s.input[s.pos]
	s.pos++
	return ch, true
}

func (s *Scanner) Rollback() {
	if s.pos > 0 {
		s.pos--
	}
}

var transition = map[State][]State{
	S0: {S1, Se, Se},
	S1: {Se, S2, Se},
	S2: {Se, S2, Se},
	Se: {Se, Se, Se},
}

var accepts = []State{S2}

var token_table = map[State]string{
	S0: "invalid",
	S1: "invalid",
	S2: "register",
	Se: "invalid",
}

func (s *Scanner) NextWord() string {
	state := S0
	lexeme := ""
	s.stack = append(s.stack, Bad)
	for {
		if state == Se {
			break
		}

		c, ok := s.NextChar()
		if !ok {
			break
		}

		lexeme += string(c)
		cat := s.CharCat(c)
		s.stack = append(s.stack, state)
		state = transition[state][cat]
	}

	for {
		if slices.Contains(accepts, state) || state == Bad {
			break
		}
		sl := len(s.stack)
		state = s.stack[sl-1]
		s.stack = s.stack[:sl-1]
		if len(lexeme) > 0 {
			lexeme = lexeme[:len(lexeme)-1]
			s.Rollback()
		}
	}

	if slices.Contains(accepts, state) {
		return token_table[state]
	} else {
		return "invalid"
	}

}

func main() {
	tests := []string{
		"r1",
		"r999",
		"r42abc",
		"abc",
		"r",
		"r0",
	}

	for _, t := range tests {
		s := new_scanner(t)

		token := s.NextWord()
		fmt.Printf("%q => MATCH %q\n", t, token)
	}
}
