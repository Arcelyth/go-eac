package main

import (
	"fmt"
)

type State byte

const (
	S0 State = iota
	S1
	S2
	Se
	Bad
)

type Scanner struct {
	input string
	pos   int
}

func newScanner(input string) *Scanner {
	return &Scanner{input: input}
}

func (s *Scanner) next() byte {
	if s.pos >= len(s.input) {
		return 0
	}
	ch := s.input[s.pos]
	s.pos++
	return ch
}

func (s *Scanner) rollback() {
	if s.pos > 0 {
		s.pos--
	}
}

func (s *Scanner) NextWord() string {
	start := s.pos
	state := S0

	for {
		ch := s.next()
		if ch == 0 {
			break
		}

		switch state {
		case 0:
			if ch == 'r' {
				state = S1
			} else {
				return "invalid"
			}
		case 1:
			if ch >= '0' && ch <= '9' {
				state = S2
			} else {
				return "invalid"
			}

		case 2:
			if ch >= '0' && ch <= '9' {
				state = S2
			} else {
				s.rollback()
				return "register"
			}
		}
	}

	if state == 2 {
		return "register"
	}

	s.pos = start
	return "invalid"
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
		fmt.Printf("%q => %s\n", t, s.NextWord())
	}
}
