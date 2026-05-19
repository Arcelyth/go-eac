package main

import (
	"fmt"
	"slices"
	"unicode"
)

type TokenType byte

type Token struct {
	value string
	kind  string
}

type Lexer struct {
	src []rune
	pos int
}

func newLexer(src string) Lexer {
	return Lexer{
		src: []rune(src),
		pos: 0,
	}
}

func (l *Lexer) NextToken() Token {
	for l.pos < len(l.src) && unicode.IsSpace(l.src[l.pos]) {
		l.pos++
	}

	if l.pos >= len(l.src) {
		return Token{kind: "eof", value: "eof"}
	}

	ch := l.src[l.pos]
	l.pos++

	switch ch {
	case '+':
		return Token{kind: "+", value: "+"}
	case '-':
		return Token{kind: "-", value: "-"}
	case '*', '×', 'x':
		return Token{kind: "×", value: string(ch)}
	case '/':
		return Token{kind: "/", value: string(ch)}
	case '(':
		return Token{kind: "(", value: "("}
	case ')':
		return Token{kind: ")", value: ")"}
	}

	if unicode.IsDigit(ch) {
		start := l.pos - 1
		for l.pos < len(l.src) && unicode.IsDigit(l.src[l.pos]) {
			l.pos++
		}
		return Token{kind: "num", value: string(l.src[start:l.pos])}
	}

	if unicode.IsLetter(ch) {
		start := l.pos - 1
		for l.pos < len(l.src) && unicode.IsDigit(l.src[l.pos]) || unicode.IsDigit(l.src[l.pos]) {
			l.pos++
		}
		return Token{kind: "name", value: string(l.src[start:l.pos])}
	}
	return Token{kind: "error", value: string(ch)}
}

type Production struct {
	id  int
	lhs string
	rhs []string
}

type Set map[string]bool

func NewSet() Set { return make(Set) }

func (s Set) Union(s2 Set) bool {
	changed := false
	for k := range s2 {
		if !s[k] {
			s[k] = true
			changed = true
		}
	}

	return changed
}

type LL1Builder struct {
	productions  []Production
	terminals    []string
	nonTerminals []string
	startSymbol  string

	firstSets  map[string]Set
	followSets map[string]Set
	table      map[string]map[string]int
}

func NewLL1Builder() *LL1Builder {
	return &LL1Builder{
		firstSets:  make(map[string]Set),
		followSets: make(map[string]Set),
		table:      make(map[string]map[string]int),
	}
}

func (b *LL1Builder) isTerminal(sym string) bool {
	if sym == "null" || sym == "eof" {
		return true
	}
	if slices.Contains(b.terminals, sym) {
		return true
	}
	return false
}

func (b *LL1Builder) GetFirst(rhs []string) Set {
	res := NewSet()
	if len(rhs) == 0 {
		res["null"] = true
		return res
	}
	for i, sym := range rhs {
		if b.isTerminal(sym) {
			res[sym] = true
			break
		}
		fSet := b.firstSets[sym]
		hasNull := false
		for k := range fSet {
			if k != "null" {
				res[k] = true
			} else {
				hasNull = true
			}
		}

		if !hasNull {
			break
		}
		// if the last one is epsilon use follow set
		if i == len(rhs)-1 {
			res["null"] = true
		}
	}
	return res
}

func (b *LL1Builder) BuildTable() {
	for _, nt := range b.nonTerminals {
		b.firstSets[nt] = NewSet()
		b.followSets[nt] = NewSet()
		b.table[nt] = make(map[string]int)
		for _, t := range b.terminals {
			b.table[nt][t] = -1
		}
		b.table[nt]["eof"] = -1
	}

	// First
	for {
		changed := false
		for _, p := range b.productions {
			fStr := b.GetFirst(p.rhs)
			if b.firstSets[p.lhs].Union(fStr) {
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	// Follow
	b.followSets[b.startSymbol]["eof"] = true
	for {
		changed := false
		for _, p := range b.productions {
			for i, sym := range p.rhs {
				if b.isTerminal(sym) {
					continue
				}
				trailer := p.rhs[i+1:]
				fStr := b.GetFirst(trailer)
				for k := range fStr {
					if k != "null" {
						if !b.followSets[sym][k] {
							b.followSets[sym][k] = true
							changed = true
						}
					}
				}
				if fStr["null"] || i == len(p.rhs)-1 {
					if b.followSets[sym].Union(b.followSets[p.lhs]) {
						changed = true
					}
				}
			}
		}
		if !changed {
			break
		}
	}

	// First+
	for _, p := range b.productions {
		fStr := b.GetFirst(p.rhs)
		firstPlus := NewSet()
		if !fStr["null"] {
			firstPlus.Union(fStr)
		} else {
			for k := range fStr {
				if k != "null" {
					firstPlus[k] = true
				}
			}
			firstPlus.Union(b.followSets[p.lhs])
		}
		for w := range firstPlus {
			if w != "null" {
				b.table[p.lhs][w] = p.id
			}
		}
	}
}

type LL1Parser struct {
	builder *LL1Builder
}

func NewLL1Parser(builder *LL1Builder) *LL1Parser {
	return &LL1Parser{builder: builder}
}

func (p *LL1Parser) Parse(input string) {
	lexer := newLexer(input)

	// word ← NextWord( );
	tok := lexer.NextToken()
	word := tok.kind

	var stack []string

	// push eof onto Stack;
	stack = append(stack, "eof")
	// push the start symbol, S, onto Stack;
	stack = append(stack, p.builder.startSymbol)

	fmt.Printf("\n[Expression]: \"%s\"\n", input)
	fmt.Println("--------------------------------------------------------------------------------")

	// loop forever;
	for {
		if len(stack) == 0 {
			fmt.Println("Error: Stack underflow.")
			return
		}

		// focus ← top of Stack;
		focus := stack[len(stack)-1]

		fmt.Printf("stack: %-30v | focus: %-8s | token: %s\n", fmt.Sprintf("%v", stack), focus, tok.value)

		// if (focus = eof and word = eof)
		if focus == "eof" && word == "eof" {
			// then report success and exit the loop;
			fmt.Println("[Success]")
			return
		}

		// else if (focus ∈ T or focus = eof)
		if p.builder.isTerminal(focus) || focus == "eof" {
			// if focus matches word
			if focus == word {
				// pop Stack;
				stack = stack[:len(stack)-1]
				// word ← NextWord( );
				tok = lexer.NextToken()
				word = tok.kind
			} else {
				// else report an error looking for symbol at top of stack;
				fmt.Printf("\n[Error]: Expect '%s', but get '%s' (value: '%s')\n", focus, word, tok.value)
				return
			}
		} else {
			// else begin; /* focus is a nonterminal */
			row := p.builder.table[focus]
			prodID, exists := row[word]

			// if Table[focus,word] is A → B1B2 ··· Bk
			if exists && prodID != -1 {
				// pop Stack;
				stack = stack[:len(stack)-1]
				production := p.builder.productions[prodID]

				// for i ← k to 1 by -1 do;
				for i := len(production.rhs) - 1; i >= 0; i-- {
					bi := production.rhs[i]
					// if (Bi != null)
					if bi != "" && bi != "null" {
						// then push Bi onto Stack;
						stack = append(stack, bi)
					}
				}
			} else {
				// else report an error expanding focus;
				fmt.Printf("\n[Error] Empty \n")
				return
			}
		}
	}
}

func main() {
	builder := NewLL1Builder()

	builder.startSymbol = "Goal"
	builder.nonTerminals = []string{"Goal", "Expr", "Expr'", "Term", "Term'", "Factor"}
	builder.terminals = []string{"+", "-", "×", "/", "(", ")", "name", "num"}

	builder.productions = []Production{
		0:  {id: 0, lhs: "Goal", rhs: []string{"Expr"}},
		1:  {id: 1, lhs: "Expr", rhs: []string{"Term", "Expr'"}},
		2:  {id: 2, lhs: "Expr'", rhs: []string{"+", "Term", "Expr'"}},
		3:  {id: 3, lhs: "Expr'", rhs: []string{"-", "Term", "Expr'"}},
		4:  {id: 4, lhs: "Expr'", rhs: []string{}}, // rhs empty == null
		5:  {id: 5, lhs: "Term", rhs: []string{"Factor", "Term'"}},
		6:  {id: 6, lhs: "Term'", rhs: []string{"×", "Factor", "Term'"}},
		7:  {id: 7, lhs: "Term'", rhs: []string{"/", "Factor", "Term'"}},
		8:  {id: 8, lhs: "Term'", rhs: []string{}},
		9:  {id: 9, lhs: "Factor", rhs: []string{"(", "Expr", ")"}},
		10: {id: 10, lhs: "Factor", rhs: []string{"num"}},
		11: {id: 11, lhs: "Factor", rhs: []string{"name"}},
	}

	builder.BuildTable()

	fmt.Println("================================================================================")
	fmt.Println("                                LL(1) Table                                     ")
	fmt.Println("================================================================================")
	cols := append(builder.terminals, "eof")
	fmt.Printf("%-10s", "NonTerminal")
	for _, c := range cols {
		fmt.Printf("%-7s", c)
	}
	fmt.Println()
	for _, nt := range builder.nonTerminals {
		fmt.Printf("%-10s", nt)
		for _, c := range cols {
			val := builder.table[nt][c]
			if val == -1 {
				fmt.Printf("%-7s", "—")
			} else {
				fmt.Printf("%-7d", val)
			}
		}
		fmt.Println()
	}
	fmt.Println("================================================================================")
	parser := NewLL1Parser(builder)

	parser.Parse("1 + 4 x 5")
	parser.Parse("(varName - 100) / 2")
	parser.Parse("1 + + 2")

}
