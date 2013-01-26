package vm

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
