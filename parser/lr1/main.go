package main

import (
	"fmt"
	"os"
	"strconv"
)

type ActionKind string

const (
	SHIFT  ActionKind = "shift"
	REDUCE ActionKind = "reduce"
	ACCEPT ActionKind = "accept"
	ERROR  ActionKind = "error"
)

type Action struct {
	kind ActionKind
	// shift -> target; reduce  -> production
	target int
}

type Production struct {
	id  int
	lhs string
	rhs []string
}

// 0: Goal -> Expr
// 1: Expr -> ( Expr )
// 2: Expr -> ( )
var productions = []Production{
	0: {id: 0, lhs: "Goal", rhs: []string{"Expr"}},
	1: {id: 1, lhs: "Expr", rhs: []string{"(", "Expr", ")"}},
	2: {id: 2, lhs: "Expr", rhs: []string{"(", ")"}},
}

var ActionTable = map[int]map[string]Action{
	0: {"(": {kind: SHIFT, target: 2}},
	1: {"eof": {kind: ACCEPT}},
	2: {"(": {kind: SHIFT, target: 2}, ")": {kind: SHIFT, target: 4}},
	3: {")": {kind: SHIFT, target: 5}},
	4: {")": {kind: REDUCE, target: 2}, "eof": {kind: REDUCE, target: 2}},
	5: {")": {kind: REDUCE, target: 1}, "eof": {kind: REDUCE, target: 1}},
}

var GotoTable = map[int]map[string]int{
	0: {"Expr": 1},
	2: {"Expr": 3},
}

func Fail(message string) {
	fmt.Printf("\n[Error]: %s\n", message)
	os.Exit(1)
}

func main() {
	tokens := []string{"(", "(", ")", ")", "eof"}
	tokenIdx := 0

	nextWord := func() string {
		if tokenIdx >= len(tokens) {
			return "eof"
		}
		t := tokens[tokenIdx]
		tokenIdx++
		return t
	}

	var stack []string

	// push $;
	stack = append(stack, "$")

	// push start state, s0;
	stack = append(stack, "0")

	// word ← NextWord( );
	word := nextWord()

	fmt.Println("================================================================================")
	fmt.Println("                                   LR(1)                                        ")
	fmt.Println("================================================================================")

	// while (true) do;
	for {
		// state ← top of stack;
		topElement := stack[len(stack)-1]
		state, err := strconv.Atoi(topElement)
		if err != nil {
			Fail(fmt.Sprintf("Top element %s has wrong type", topElement))
		}

		fmt.Printf("stack: %-40v | word: %s\n", fmt.Sprintf("%v", stack), word)

		row, hasState := ActionTable[state]
		if !hasState {
			Fail(fmt.Sprintf("Status %d has empty line in Action table", state))
		}
		action, hasAction := row[word]
		if !hasAction {
			Fail(fmt.Sprintf("\nStatus %d cannot accept '%s'", state, word))
		}

		// if Action[state,word] = ‘‘reduce A→β’’ then begin;
		if action.kind == REDUCE {
			prod := productions[action.target]
			betaLen := len(prod.rhs)

			fmt.Printf(" => [Action]: Reduce -> Production %d: %s -> %v\n", prod.id, prod.lhs, prod.rhs)

			// pop 2 × | β | symbols;
			popCount := 2 * betaLen
			if len(stack) < popCount {
				Fail("Length of stack is not enough")
			}
			stack = stack[:len(stack)-popCount]

			// state ← top of stack;
			underStateStr := stack[len(stack)-1]
			underState, _ := strconv.Atoi(underStateStr)

			// push A;
			stack = append(stack, prod.lhs)

			// push Goto[state, A];
			gotoRow, hasGotoRow := GotoTable[underState]
			nextStateID, hasGoto := gotoRow[prod.lhs]
			if !hasGotoRow || !hasGoto {
				Fail(fmt.Sprintf("Status %d cannot do goto throught '%s'", underState, prod.lhs))
			}
			stack = append(stack, strconv.Itoa(nextStateID))

			// end;
		} else if action.kind == SHIFT {
			// else if Action[state,word] = ‘‘shift si’’ then begin;
			targetState := action.target
			fmt.Printf(" => [Action]: Shift -> Symbol '%s' to Status %d\n", word, targetState)

			// push word;
			stack = append(stack, word)

			// push si ;
			stack = append(stack, strconv.Itoa(targetState))

			// word ← NextWord( );
			word = nextWord()

			// end;
		} else if action.kind == ACCEPT {
			// else if Action[state,word] = ‘‘accept’’ then break;
			fmt.Println("  └─> [Action]: Accept")
			break
		} else {
			// else Fail( );
			Fail("Unknown error")
		}
	}

	// report success;
	fmt.Println("[Success]")
}
