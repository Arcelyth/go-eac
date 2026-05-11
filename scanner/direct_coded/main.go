// just replace the transition table with switch

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

type Scanner struct {
	input  string
	pos    int
	failed map[State]map[int]bool
}

func newScanner(input string) Scanner {
	failed := make(map[State]map[int]bool)
	for _, st := range []State{S0, S1, S2, Se, Bad} {
		failed[st] = make(map[int]bool)
	}
	return Scanner{
		input:  input,
		pos:    0,
		failed: failed,
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
		stack = append(stack, Entry{state, s.pos})
		switch state {
		case S0:
			if c == 'r' {
				state = S1
			} else {
				state = Se
			}

		case S1:
			if c >= '0' && c <= '9' {
				state = S2
			} else {
				state = Se
			}

		case S2:
			if c >= '0' && c <= '9' {
				state = S2
			} else {
				state = Se
			}
		}
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
		s := newScanner(t)

		token := s.NextWord()
		fmt.Printf("%q => MATCH %q\n", t, token)
	}
}
