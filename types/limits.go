package types

import (
	"math"
)

// Size of opcode arguments
const (
	sizeC  = 9
	sizeB  = 9
	sizeBx = sizeC + sizeB
	sizeA  = 8
	sizeAx = sizeC + sizeB + sizeA
	sizeOp = 6
)

// Position of opcode arguments
const (
	posOp = 0
	posA  = posOp + sizeOp
	posC  = posA + sizeA
	posB  = posC + sizeC
	posBx = posC
	posAx = posA
)

// INT_MAX - 2 "for safety"
// https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/a9PitPAHSSU
const (
	MAX_INT    = math.MaxInt32 - 2 // TODO : Int32 or int64? Check Lua on Linux 64bit
	MAXARG_Ax  = ((1 << sizeAx) - 1)
	MAXARG_Bx  = ((1 << sizeBx) - 1)
	MAXARG_sBx = (MAXARG_Bx >> 1)
	MAXARG_A   = ((1 << sizeA) - 1)
	MAXARG_B   = ((1 << sizeB) - 1)
	MAXARG_C   = ((1 << sizeC) - 1)

	MAXINDEXRK = (BITRK - 1)
)
