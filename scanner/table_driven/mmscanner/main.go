// maximal munch scanner,

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

type Entry struct {
	state State
	pos   int
}

type MMScanner struct {
	input  string
	pos    int
	failed map[State]map[int]bool
}

func new_scanner(input string) MMScanner {
	failed := make(map[State]map[int]bool)
	for _, st := range []State{S0, S1, S2, Se, Bad} {
		failed[st] = make(map[int]bool)
	}
	return MMScanner{
		input:  input,
		pos:    0,
		failed: failed,
	}
}

func (s *MMScanner) CharCat(ch byte) charType {
	switch {
	case ch == 'r':
		return Register
	case ch >= '0' && ch <= '9':
		return Data
	default:
		return Other
	}
}

func (s *MMScanner) NextChar() (byte, bool) {
	if s.pos >= len(s.input) {
		return 0, false
	}
	ch := s.input[s.pos]
	s.pos++
	return ch, true
}

func (s *MMScanner) Rollback() {
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

func (s *MMScanner) NextWord() string {
	state := S0
	lexeme := ""
	stack := []Entry{{Bad, -1}}
	for {
		if state == Se {
			break
		}
		c, ok := s.NextChar()
		if !ok {
			break
		}

		if s.failed[state][s.pos] {
			break
		}

		if slices.Contains(accepts, state) {
			stack = []Entry{}
		}

		lexeme += string(c)
		cat := s.CharCat(c)
		stack = append(stack, Entry{state, s.pos})
		state = transition[state][cat]
	}

	for {
		if slices.Contains(accepts, state) || state == Bad {
			break
		}
		s.failed[state][s.pos] = true
		sl := len(stack)
		state = stack[sl-1].state
		stack = stack[:sl-1]
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
