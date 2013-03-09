package types

import (
	"fmt"
)

// Holds pointers to values (pointer to empty interface - yes, I know, but it is 
// required because the interface may hold the value inline for numbers and bools):
// http://play.golang.org/p/e2Ptu8puSZ

const (
	_INITIAL_STACK_CAP  = 10
	_INITIAL_UPVALS_CAP = 5
)

type State struct {
	Stack       []Value
	Top         int // index of the first free slot in the stack
	Globals     Table
	CI          *CallInfo
	OpenUpVals  []int    // TODO : index into the stack or struct?
	OpCodeDebug []OpCode // TODO : Very temporary, find a better solution for testing... hooks?
}

func NewState(entryPoint *Prototype) *State {
	s := &State{
		Stack:      make([]Value, _INITIAL_STACK_CAP),
		Globals:    NewTable(),
		OpenUpVals: make([]int, _INITIAL_UPVALS_CAP),
	}

	cl := NewClosure(entryPoint)
	if l := len(entryPoint.Upvalues); l == 1 {
		// 1 upvalue = globals table as upvalue
		v := Value(s.Globals)
		cl.UpVals[0] = v
	} else if l > 1 {
		panic(fmt.Sprintf("too many upvalues expected for entry point: %d", l))
	}

	// Push the closure on the stack
	s.CheckStack(cl.P.Meta.MaxStackSize + 1) // +1 for the closure itself
	s.Stack[s.Top] = cl
	s.Top++
	return s
}

func (s *State) CheckStack(needed byte) {
	oriAdr := &s.Stack[0]

	missing := (s.Top + int(needed)) - cap(s.Stack)
	for i := 0; i < missing; i++ {
		s.Stack = append(s.Stack, nil)
	}

	if oriAdr != &s.Stack[0] {
		// Re-capture the frame for all existing CallInfos
		for ci := s.CI; ci != nil; ci = ci.Prev {
			ci.captureFrame(s)
		}
	}
}

// TODO : Very very temporary...
func (s *State) Dump() {
	fmt.Print(">>")
	for i, v := range s.Stack {
		if i == s.Top {
			fmt.Print("\u01AE")
		}
		if i == s.CI.FuncIndex {
			fmt.Print("\u0192")
		}
		if i == s.CI.Base {
			fmt.Print("\u0251")
		}
		fmt.Printf("%v, ", v)
	}
	for j := len(s.Stack); j <= s.Top; j++ {
		if j == s.Top {
			fmt.Print("\u01AE")
		}
		if j == s.CI.FuncIndex {
			fmt.Print("\u0192")
		}
		if j == s.CI.Base {
			fmt.Print("\u0251")
		}
		fmt.Println(", ", j)
	}
	fmt.Print("<<\n")
}

type CallInfo struct {
	Frame      []Value
	Cl         *Closure
	FuncIndex  int
	NumResults int
	CallStatus byte
	PC         int
	Base       int
	Prev       *CallInfo
}

func (s *State) NewCallInfo(cl *Closure, idx int, nRets int) {
	// Make sure the stack has enough slots
	s.CheckStack(cl.P.Meta.MaxStackSize)

	// Complete the arguments
	n := s.Top - idx - 1
	for ; n < int(cl.P.Meta.NumParams); n++ {
		s.Stack[s.Top] = nil
		s.Top++
	}

	var base int
	if cl.P.Meta.IsVarArg == 0 {
		base = idx + 1
	} else {
		if !(n >= int(cl.P.Meta.NumParams)) {
			panic(fmt.Sprintf("expected actual number of args (%d) to be greater than or equal to the fixed number of args (%d)", n, cl.P.Meta.NumParams))
		}
		fixed := s.Top - n
		base = s.Top
		for i := 0; i < int(cl.P.Meta.NumParams); i++ {
			s.Stack[s.Top] = s.Stack[fixed+i]
			s.Top++
			s.Stack[fixed+i] = nil
		}
	}

	ci := new(CallInfo)
	ci.Cl = cl
	ci.FuncIndex = idx
	ci.NumResults = nRets
	ci.CallStatus = 0 // TODO : For now, ignore
	ci.PC = 0
	ci.Base = base
	ci.Prev = s.CI
	s.Top = base + int(cl.P.Meta.MaxStackSize)
	ci.captureFrame(s)

	s.CI = ci
}

// Because checkStack() can reallocate a new array for the stack, the frame
// may become invalid. This gets called when required to make sure that the
// frame slice always points to the stack array.
func (ci *CallInfo) captureFrame(s *State) {
	ci.Frame = s.Stack[ci.Base:(ci.Base + int(ci.Cl.P.Meta.MaxStackSize))]
}
