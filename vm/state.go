package vm

import (
	"fmt"
	"github.com/PuerkitoBio/lune/types"
)

// Holds pointers to values (pointer to empty interface - yes, I know, but it is 
// required because the interface may hold the value inline for numbers and bools):
// http://play.golang.org/p/e2Ptu8puSZ
type Stack struct {
	top int // First free slot
	stk []*types.Value
}

func newStack() *Stack {
	return new(Stack)
}

type State struct {
	Stack   *Stack
	Globals types.Table
}

func (s *Stack) Get(idx int) *types.Value {
	return s.stk[idx]
}

func (s *Stack) push(v types.Value) {
	s.stk[s.top] = &v
	s.top++
}

func (s *Stack) checkStack(needed byte) {
	missing := (s.top + int(needed) - 1) - cap(s.stk)
	for i := 0; i < missing; i++ {
		var v types.Value
		v = nil
		s.stk = append(s.stk, &v)
	}
}

func (s *Stack) dumpStack() {
	fmt.Println("*** DUMP STACK ***")
	for i, v := range s.stk {
		if v == nil {
			fmt.Println(i, v)
		} else if f, ok := (*v).(*types.Closure); ok {
			fmt.Println(i, f.P.Source, f.P.Meta.LineDefined)
		} else {
			fmt.Println(i, *v)
		}
	}
}

func NewState(entryPoint *types.Prototype) *State {
	s := &State{newStack(), make(types.Table)}

	cl := types.NewClosure(entryPoint)
	if l := len(entryPoint.Upvalues); l == 1 {
		// 1 upvalue = globals table as upvalue
		v := types.Value(s.Globals)
		cl.UpVals[0] = &v
	} else if l > 1 {
		// TODO : panic?
		panic("too many upvalues expected for entry point")
	}

	// Push the closure on the stack
	s.Stack.checkStack(cl.P.Meta.MaxStackSize)
	s.Stack.push(cl)
	return s
}

type CallInfo struct {
	Cl         *types.Closure
	FuncIndex  int
	NumResults int
	CallStatus byte
	PC         int
	Base       int
	Prev       *CallInfo
}

func newCallInfo(s *State, fIdx int, prev *CallInfo) *CallInfo {
	// Get the function's closure at this stack index
	f := s.Stack.Get(fIdx)
	cl := (*f).(*types.Closure)

	// Make sure the stack has enough slots
	s.Stack.checkStack(cl.P.Meta.MaxStackSize)

	// Complete the arguments
	n := s.Stack.top - fIdx - 1
	for ; n < int(cl.P.Meta.NumParams); n++ {
		s.Stack.push(nil)
	}

	ci := new(CallInfo)
	ci.Cl = cl
	ci.FuncIndex = fIdx
	ci.NumResults = 0 // TODO : For now, ignore, someday will be passed
	ci.CallStatus = 0 // TODO : For now, ignore
	ci.PC = 0
	ci.Base = fIdx + 1 // TODO : For now, considre the base to be fIdx + 1, will have to manage varargs someday
	ci.Prev = prev

	return ci
}

/*
type gState struct {
	strt strTable
}

func newGState() *gState {
	return &gState{newStrTable()}
}

type lState struct {
	g      *gState
	stk    *stack
	ci     *callInfo
	baseCi callInfo
}

func NewState() *lState {
	return &lState{newGState(), newStack(), nil, callInfo{}}
}

type callInfo struct {
	funcStkIdx uint
	prev, next *callInfo
	nResults   uint8
}
*/
