package vm

import (
	"github.com/PuerkitoBio/lune/types"
)

// Holds pointers to values (pointer to empty interface - yes, I know, but it is 
// required because the interface may hold the value inline for numbers and bools):
// http://play.golang.org/p/e2Ptu8puSZ
type Stack struct {
	top int // First free slot
	stk []*types.Value
}

type State struct {
	stack   *Stack
	globals types.Table
}

func (s *State) Get(idx int) *types.Value {
	return s.stack.stk[idx]
}

func (s *State) push(v types.Value) {
	s.stack.stk[s.stack.top] = &v
	s.stack.top++
}

func (s *State) checkStack(needed byte) {
	missing := cap(s.stack.stk) - (s.stack.top + int(needed) + 1) // i.e. cap=10, top=7 and is last used - so 8 slots taken, needed=3: 10-(7 + 3 + 1)
	if missing > 0 {
		dummy := make([]*types.Value, missing)
		s.stack.stk = append(s.stack.stk, dummy...)
	}
}

func NewState(entryPoint *types.Prototype) *State {
	s := &State{new(Stack)}
	s.globals = new(types.Table)

	if l := len(entryPoint.Upvalues); l == 1 {
		// 1 upvalue = globals table as upvalue
	} else if l > 1 {
		// TODO : panic?
		panic("too many upvalues expected for entry point")
	}

	// Push on the stack
	s.push(entryPoint)
	return s
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
