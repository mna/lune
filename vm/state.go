package vm

import (
	"github.com/PuerkitoBio/lune/serializer"
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
	stack *Stack
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

func NewState(entryPoint *serializer.Prototype) *State {
	s := &State{new(Stack)}

	// TODO : For now, Index 0 is always the Global's table, probably more complex than this if the index is relative to base of the function
	s.push(new(Table))
	// Index 1
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
