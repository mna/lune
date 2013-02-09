package types

import (
	"fmt"
)

// Holds pointers to values (pointer to empty interface - yes, I know, but it is 
// required because the interface may hold the value inline for numbers and bools):
// http://play.golang.org/p/e2Ptu8puSZ

type Stack struct {
	Top int // First free slot
	stk []Value
}

func newStack() *Stack {
	return new(Stack)
}

type State struct {
	Stack   *Stack
	Globals Table
	CI      *CallInfo
}

func (s *Stack) Get(idx int) Value {
	return s.stk[idx]
}

func (s *Stack) Push(v Value) {
	s.stk[s.Top] = v
	s.Top++
}

func (s *Stack) checkStack(needed byte) {
	missing := (s.Top + int(needed) + 1) - cap(s.stk)
	for i := 0; i < missing; i++ {
		var v Value
		v = nil
		s.stk = append(s.stk, v)
	}
}

func (s *Stack) DumpStack() {
	fmt.Println("*** DUMP STACK ***")
	for i, v := range s.stk {
		if f, ok := v.(*Closure); ok {
			fmt.Println(i, f.P.Source, f.P.Meta.LineDefined)
		} else {
			fmt.Println(i, v)
		}
	}
}

func NewState(entryPoint *Prototype) *State {
	s := &State{
		Stack:   newStack(),
		Globals: NewTable(),
	}

	cl := NewClosure(entryPoint)
	if l := len(entryPoint.Upvalues); l == 1 {
		// 1 upvalue = globals table as upvalue
		v := Value(s.Globals)
		cl.UpVals[0] = v
	} else if l > 1 {
		// TODO : panic?
		panic("too many upvalues expected for entry point")
	}

	// Push the closure on the stack
	s.Stack.checkStack(cl.P.Meta.MaxStackSize)
	s.Stack.Push(cl)
	return s
}

func (s *State) NewCallInfo(fIdx int, prev *CallInfo) {
	// Get the function's closure at this stack index
	f := s.Stack.Get(fIdx)
	cl := f.(*Closure)

	// Make sure the stack has enough slots
	s.Stack.checkStack(cl.P.Meta.MaxStackSize)

	// Complete the arguments
	n := s.Stack.Top - fIdx - 1
	for ; n < int(cl.P.Meta.NumParams); n++ {
		s.Stack.Push(nil)
	}

	ci := new(CallInfo)
	ci.Cl = cl
	ci.FuncIndex = fIdx
	ci.NumResults = 0 // TODO : For now, ignore, someday will be passed
	ci.CallStatus = 0 // TODO : For now, ignore
	ci.PC = 0
	ci.Base = fIdx + 1 // TODO : For now, considre the base to be fIdx + 1, will have to manage varargs someday
	ci.Prev = prev
	ci.Frame = s.Stack.stk[ci.Base:]

	s.CI = ci
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
